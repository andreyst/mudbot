package server

import (
	"encoding/json"
	"github.com/gorilla/websocket"
)

func (s *Server) sendAtlasData(c *websocket.Conn, event string) error {
	atlasData := s.atlasDataProvider()
	message := struct {
		Data  interface{}
		Event string
	}{
		Data:  atlasData,
		Event: event,
	}

	messageJson, messageMarshalErr := json.MarshalIndent(message, "", "  ")
	if messageMarshalErr != nil {
		return messageMarshalErr
	}

	writeMessageErr := c.WriteMessage(websocket.TextMessage, messageJson)
	if writeMessageErr != nil {
		return writeMessageErr
	}

	return nil
}

func (s *Server) wsUpdater(conn *websocket.Conn, updateCh chan string, closeCh chan interface{}) {
	defer close(closeCh)

	sendRoomsErr := s.sendAtlasData(conn, "init")
	if sendRoomsErr != nil {
		s.logger.Errorf("Initial send rooms err: %v", sendRoomsErr)
		return
	}

	for {
		var event string
		var ok bool
		select {
		case event, ok = <-updateCh:
			if !ok {
				return
			}
		}

		sendAtlasDataErr := s.sendAtlasData(conn, event)
		if sendAtlasDataErr != nil {
			s.logger.Debugf("send atlas Data err: %v", sendAtlasDataErr)
			return
		}
	}
}
