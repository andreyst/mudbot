package atlas

import "encoding/json"

func (a *Atlas) Parse(message []byte) {
	var data interface{}

	json.Unmarshal(message, &data)
}
