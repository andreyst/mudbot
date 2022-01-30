package server

import (
	"encoding/json"
	"github.com/gorilla/websocket"
)

func (s *Server) sendAtlasData(c *websocket.Conn) error {
	atlasData := s.atlasDataProvider()

	messageJson, messageMarshalErr := json.MarshalIndent(atlasData, "", "  ")
	if messageMarshalErr != nil {
		return messageMarshalErr
	}

	writeMessageErr := c.WriteMessage(websocket.TextMessage, messageJson)
	if writeMessageErr != nil {
		return writeMessageErr
	}

	return nil
}

func (s *Server) wsUpdater(conn *websocket.Conn, updateCh, closeCh chan interface{}) {
	defer close(closeCh)

	sendRoomsErr := s.sendAtlasData(conn)
	if sendRoomsErr != nil {
		s.logger.Errorf("Initial send rooms err: %v", sendRoomsErr)
		return
	}

	for {
		select {
		case _, ok := <-updateCh:
			if !ok {
				return
			}
		}

		sendAtlasDataErr := s.sendAtlasData(conn)
		if sendAtlasDataErr != nil {
			s.logger.Debugf("send atlas data err: %v", sendAtlasDataErr)
			return
		}
	}
}
