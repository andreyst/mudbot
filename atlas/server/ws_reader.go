package server

import (
	"encoding/json"
	"github.com/gorilla/websocket"
)

type Message struct {
	ShiftRoomCommand   *ShiftRoomCommand   `json:",omitempty"`
	DeleteRoomCommand  *DeleteRoomCommand  `json:",omitempty"`
	LinkRoomsCommand   *LinkRoomsCommand   `json:",omitempty"`
	UnlinkRoomsCommand *UnlinkRoomsCommand `json:",omitempty"`
}

type ShiftRoomCommand struct {
	RoomId    int
	Direction string
}

type DeleteRoomCommand struct {
	RoomId int
}

type LinkRoomsCommand struct {
	FromRoomId    int
	DirectionFrom string
	ToRoomId      int
	DirectionTo   string
}

type UnlinkRoomsCommand struct {
	FromRoomId int
	ToRoomId   int
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

	if message.ShiftRoomCommand != nil {
		s.OnShiftRoom(*message.ShiftRoomCommand)
	} else if message.DeleteRoomCommand != nil {
		s.OnDeleteRoom(*message.DeleteRoomCommand)
	} else if message.LinkRoomsCommand != nil {
		s.OnLinkRooms(*message.LinkRoomsCommand)
	} else if message.UnlinkRoomsCommand != nil {
		s.OnUnlinkRooms(*message.UnlinkRoomsCommand)
	} else {
		s.logger.Infof("Unknown command from client: %v", message)
	}
}
