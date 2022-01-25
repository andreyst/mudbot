package atlas

import (
	"encoding/json"
	"strings"
)

func (a Atlas) Html() string {
	tpl := `
<html>
<head>
	<style>
	</style>
</head>
<body>
<div id="map"></div>
<script>
var rooms = Object.values({roomsJs})

for (let room of rooms) {
	console.log(room)
}
</script>
</body>
</html>
`

	//var roomsJs strings.Builder
	roomsJsBytes, marshalErr := json.MarshalIndent(a.Rooms, "", "  ")
	if marshalErr != nil {
		a.logger.Errorf("cannot marshal rooms json: %v", marshalErr)
	}
	//for roomId, room := range a.Rooms {
	//	roomIdStr := strconv.FormatInt(roomId, 10)
	//	roomsJs.WriteString("  r" + roomIdStr + " [label=\"" + strings.ReplaceAll(room.Name, "\"", "\\\"") + "\"]\n")
	//	for exitDir, exitRoomId := range room.Exits {
	//		exitRoomIdStr := strconv.FormatInt(exitRoomId, 10)
	//		roomsJs.WriteString("  r" + roomIdStr + " -> r" + exitRoomIdStr + " [ label=\"" + exitDir.String() + "\" ];\n")
	//	}
	//	roomsJs.WriteString("\n")
	//}

	return strings.ReplaceAll(tpl, "{roomsJs}", string(roomsJsBytes))
}
