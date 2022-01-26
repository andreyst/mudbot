package atlas

import _ "embed"

import (
	"encoding/json"
	"strings"
)

//go:embed atlas_tpl.html
var htmlTpl string

func (a Atlas) Html() string {
	// TODO: Add reading from file if exists
	roomsJsBytes, marshalErr := json.MarshalIndent(a.Rooms, "", "  ")
	if marshalErr != nil {
		a.logger.Errorf("cannot marshal rooms json: %v", marshalErr)
	}

	return strings.ReplaceAll(htmlTpl, "{roomsJs}", string(roomsJsBytes))
}
