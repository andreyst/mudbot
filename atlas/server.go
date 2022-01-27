package atlas

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"mudbot/botutil"
	"net/http"
	"os"
	"sync"
)

var upgrader = websocket.Upgrader{} // use default options

type DataProvider func() (map[int64]Room, Coordinates)

type Server struct {
	dataProvider DataProvider

	chNum       int
	updateChans map[int]chan interface{}
	closeChans  map[int]chan interface{}

	mu sync.Mutex

	logger *zap.SugaredLogger
}

func NewServer(dataProvider DataProvider) *Server {
	s := Server{
		dataProvider: dataProvider,
		updateChans:  make(map[int]chan interface{}),
		closeChans:   make(map[int]chan interface{}),
		logger:       botutil.NewLogger("atlas_server"),
	}

	return &s
}

func (s *Server) makeHomeHandler(tpl string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		_, writeErr := w.Write([]byte(tpl))
		if writeErr != nil {
			s.logger.Errorf("Error writing atlas html: %v\n", writeErr)
		}
	}
}

func sendRooms(dataProvider DataProvider, c *websocket.Conn) error {
	rooms, coordinates := dataProvider()

	message := struct {
		Rooms       map[int64]Room
		Coordinates Coordinates
	}{
		rooms,
		coordinates,
	}

	messageJson, messageMarshalErr := json.MarshalIndent(message, "", "  ")
	if messageMarshalErr != nil {
		return messageMarshalErr
	}

	writeMessageErr := c.WriteMessage(websocket.TextMessage, []byte(messageJson))
	if writeMessageErr != nil {
		return writeMessageErr
	}

	return nil
}

func (s *Server) makeRoomsHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		updateCh, closeCh := s.makeUpdateChan()
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			s.logger.Debugf("ws upgrade err: %v", err)
			return
		}
		defer c.Close()

		s.logger.Debugf("Ws client connected")

		sendRoomsErr := sendRooms(s.dataProvider, c)
		if sendRoomsErr != nil {
			s.logger.Debugf("send rooms err: %v", sendRoomsErr)
			return
		}

		for {
			select {
			case _, ok := <-updateCh:
				if !ok {
					return
				}
			}

			sendRoomsErr := sendRooms(s.dataProvider, c)
			if sendRoomsErr != nil {
				s.logger.Debugf("send rooms err: %v", sendRoomsErr)
				break
			}
		}
		close(closeCh)
	}
}

func (s *Server) Start(a Atlas) {
	http.HandleFunc("/", s.makeHomeHandler(a.getHtmlTemplate()))
	http.HandleFunc("/rooms", s.makeRoomsHandler())
	go func() {
		host := os.Getenv("ATLAS_SERVER_HOST")
		port := os.Getenv("ATLAS_SERVER_PORT")
		addr := fmt.Sprintf("%s:%s", host, port)
		a.logger.Infof("Starting atlas server at %s", addr)
		a.logger.Fatal(http.ListenAndServe(addr, nil))
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

func (s *Server) sendUpdates() {
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
