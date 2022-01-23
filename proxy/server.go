package proxy

import (
	"mudbot/botutil"
	"net"
	"sync"

	"go.uber.org/zap"
)

type Server struct {
	localAddr  string
	remoteAddr string

	listener net.Listener
	stopped  bool
	workers  []*Worker

	onClientFlush AccumulatorCallback
	onMudFlush    AccumulatorCallback

	wg sync.WaitGroup

	logger *zap.SugaredLogger
}

func NewServer(localAddr string, remoteAddr string, onClientFlush AccumulatorCallback, onMudFlush AccumulatorCallback) *Server {
	s := &Server{
		localAddr:     localAddr,
		remoteAddr:    remoteAddr,
		onClientFlush: onClientFlush,
		onMudFlush:    onMudFlush,
	}
	s.logger = botutil.NewLogger("server")

	return s
}

func (s *Server) Start() {
	s.logger.Infof("Starting proxy with %v->%v", s.localAddr, s.remoteAddr)

	s.wg.Add(1)

	var listenErr error
	s.listener, listenErr = net.Listen("tcp", s.localAddr)
	if listenErr != nil {
		s.logger.Fatalf("cannot listen, exiting: %v", listenErr)
	}

	for !s.stopped {
		conn, acceptErr := s.listener.Accept()
		if acceptErr != nil {
			if s.stopped {
				break
			}
			s.logger.Fatalf("accept failed: %v", acceptErr)
		}
		go s.startWorker(conn, s.remoteAddr)
	}

	// Here and not in stop() to prevent race
	// between stop() and new worker appearing while stopping
	s.logger.Debug("stopping %v workers", len(s.workers))
	for _, worker := range s.workers {
		worker.stop()
		worker.wg.Wait()
	}

	s.logger.Infof("server shut down")
	s.wg.Done()
}

func (s *Server) Stop(wait bool) {
	s.logger.Infof("shutting down server")
	s.stopped = true
	s.listener.Close()

	if wait {
		s.wg.Wait()
	}
}

func (s *Server) startWorker(local net.Conn, remoteAddr string) {
	remote, err := net.Dial("tcp", remoteAddr)
	if remote == nil {
		s.logger.Errorf("remote dial failed: %v", err)
		return
	}

	worker := NewWorker(local, remote, s.onClientFlush, s.onMudFlush)
	s.workers = append(s.workers, worker)
	go worker.Run()
}

func (s *Server) SendToMud(str string, echo bool) {
	s.workers[0].sendToMud([]byte(str))
	if echo {
		s.workers[0].sendToClient([]byte("%bot: " + str + "\n"))
	}
}

func (s *Server) SendToClient(str string) {
	s.workers[0].sendToClient([]byte("%bot: " + str + "\n"))
}
