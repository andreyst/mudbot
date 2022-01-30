package server

import (
	"encoding/json"
	"github.com/gorilla/websocket"
)

type Message struct {
	Command       string
	ShiftCommand  *ShiftCommand  `json:",omitempty"`
	DeleteCommand *DeleteCommand `json:",omitempty"`
}

type ShiftCommand struct {
	RoomId    int
	Direction string
}

type DeleteCommand struct {
	RoomId int
}

func (s *Server) wsReader(conn *websocket.Conn) {
	defer conn.Close()

	for {
		msgType, msgBytes, readMessageErr := conn.ReadMessage()
		if readMessageErr != nil {
			s.logger.Infof("Read from client failed: %v", readMessageErr)
			return
		}

		switch msgType {
		case websocket.TextMessage:
			s.parseClientMessage(msgBytes)
		default:
			s.logger.Warnf("Unsupported message type from client: %v", msgType)
		}
	}
}

func (s *Server) parseClientMessage(msgBytes []byte) {
	s.logger.Debugf("Read message from client: %v", string(msgBytes))

	var message Message
	unmarshalErr := json.Unmarshal(msgBytes, &message)
	if unmarshalErr != nil {
		s.logger.Infof("Failed to unmarshal json from client: %v", unmarshalErr)
		return
	}

	switch message.Command {
	case "shift":
		if message.ShiftCommand == nil {
			s.logger.Infof("No command payload from client")
		} else {
			s.OnShift(*message.ShiftCommand)
		}
	case "delete":
		if message.DeleteCommand == nil {
			s.logger.Infof("No command payload from client")
		} else {
			s.OnDelete(*message.DeleteCommand)
		}
	default:
		s.logger.Infof("Unknown command from client: %v", message.Command)
	}
}
