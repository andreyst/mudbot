package server

import (
	"encoding/json"
	"github.com/gorilla/websocket"
)

type Message struct {
	Command           string
	ShiftRoomCommand  *ShiftRoomCommand  `json:",omitempty"`
	DeleteRoomCommand *DeleteRoomCommand `json:",omitempty"`
}

type ShiftRoomCommand struct {
	RoomId    int
	Direction string
}

type DeleteRoomCommand struct {
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
	case "shift_room":
		if message.ShiftRoomCommand == nil {
			s.logger.Infof("No command payload from client")
		} else {
			s.OnShiftRoom(*message.ShiftRoomCommand)
		}
	case "delete_room":
		if message.DeleteRoomCommand == nil {
			s.logger.Infof("No command payload from client")
		} else {
			s.OnDeleteRoom(*message.DeleteRoomCommand)
		}
	default:
		s.logger.Infof("Unknown command from client: %v", message.Command)
	}
}
