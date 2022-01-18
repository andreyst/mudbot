package proxy

import (
	"mudbot/botutil"
	"net"
	"sync"

	"go.uber.org/zap"
)

type server struct {
	logger *zap.SugaredLogger

	listener net.Listener
	stopped  bool

	workers []*Worker

	wg sync.WaitGroup
}

func NewServer() *server {
	s := &server{}
	s.logger = botutil.NewLogger("server")

	return s
}

func (s *server) Start(localAddr string, remoteAddr string) {
	s.wg.Add(1)

	s.logger.Infof("Starting proxy with %v->%v", localAddr, remoteAddr)
	s.logger.Debugf("Starting proxy with %v->%v", localAddr, remoteAddr)

	var listenErr error
	s.listener, listenErr = net.Listen("tcp", localAddr)
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
		go s.startWorker(conn, remoteAddr)
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

func (s *server) Stop(wait bool) {
	s.logger.Infof("shutting down server")
	s.stopped = true
	s.listener.Close()

	if wait {
		s.wg.Wait()
	}
}

func (s *server) startWorker(local net.Conn, remoteAddr string) {
	remote, err := net.Dial("tcp", remoteAddr)
	if remote == nil {
		s.logger.Errorf("remote dial failed: %v", err)
		return
	}

	worker := &Worker{}
	s.workers = append(s.workers, worker)
	go worker.Run("", local, remote)
}
