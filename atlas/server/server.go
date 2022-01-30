package server

import (
	"fmt"
	"go.uber.org/zap"
	"mudbot/botutil"
	"net/http"
	"os"
	"sync"
)

type AtlasDataProvider func() interface{}

type Server struct {
	atlasDataProvider AtlasDataProvider

	OnShift  func(ShiftCommand)
	OnDelete func(DeleteCommand)

	chNum       int
	updateChans map[int]chan interface{}
	closeChans  map[int]chan interface{}

	mu sync.Mutex

	logger *zap.SugaredLogger
}

func NewServer(atlasDataProvider AtlasDataProvider) *Server {
	s := Server{
		atlasDataProvider: atlasDataProvider,

		OnShift:  func(cmd ShiftCommand) {},
		OnDelete: func(cmd DeleteCommand) {},

		updateChans: make(map[int]chan interface{}),
		closeChans:  make(map[int]chan interface{}),

		logger: botutil.NewLogger("atlas_server"),
	}

	return &s
}

func (s *Server) Start() {
	http.HandleFunc("/", s.homeHandler)
	http.HandleFunc("/ws", s.wsHandler)
	go func() {
		host := os.Getenv("ATLAS_SERVER_HOST")
		port := os.Getenv("ATLAS_SERVER_PORT")
		addr := fmt.Sprintf("%s:%s", host, port)
		s.logger.Infof("Starting atlas server at %s", addr)
		s.logger.Fatal(http.ListenAndServe(addr, nil))
	}()
}

func (s *Server) makeUpdateChan() (chan interface{}, chan interface{}) {
	s.mu.Lock()
	s.chNum++
	updateCh := make(chan interface{})
	closeCh := make(chan interface{})
	s.updateChans[s.chNum] = updateCh
	s.closeChans[s.chNum] = closeCh
	go s.waitForCloseChan(closeCh, s.chNum)
	s.mu.Unlock()

	return updateCh, closeCh
}

func (s *Server) SendData() {
	s.mu.Lock()
	for chNum, ch := range s.updateChans {
		select {
		case ch <- 0:
			s.logger.Infof("Written update to chan %v", chNum)
		default:
			s.logger.Infof("Missed update to chan %v", chNum)
		}
	}
	s.mu.Unlock()
}

func (s *Server) waitForCloseChan(ch chan interface{}, chNum int) {
	s.logger.Debugf("Waiting for close chan %v", chNum)
	select {
	case <-ch:
		s.logger.Debugf("Closing chan %v", chNum)
		s.mu.Lock()
		delete(s.closeChans, chNum)
		delete(s.updateChans, chNum)
		s.mu.Unlock()
	}
}
