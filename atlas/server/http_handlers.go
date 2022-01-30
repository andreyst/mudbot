package server

import (
	"github.com/gorilla/websocket"
	"net/http"
)

func (s *Server) homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	_, writeErr := w.Write([]byte(s.getHtmlTemplate()))
	if writeErr != nil {
		s.logger.Errorf("Error writing atlas html: %v\n", writeErr)
	}
}

func (s *Server) wsHandler(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{}
	conn, upgradeErr := upgrader.Upgrade(w, r, nil)
	if upgradeErr != nil {
		s.logger.Debugf("ws upgrade err: %v", upgradeErr)
		return
	}
	defer conn.Close()

	s.logger.Debugf("Ws client connected")

	updateCh, closeCh := s.makeUpdateChan()

	go s.wsUpdater(conn, updateCh, closeCh)
	go s.wsReader(conn)

	<-closeCh
}
