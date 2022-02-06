package mud

import (
	"mudbot/atlas"
	"mudbot/botutil"
	"strings"
)

var prRoomExits = "\x1b[0;36m[ Exits: "
var roomNameMarker = "\x1b[1;36m"
var roomItemsMarker = "\x1b[1;33m"
var roomMobsMarker = "\x1b[1;31m"
var roomMobsEndMarker = "\x1b[0;0m"

func (p *Parser) ParseRoom(s string) (room *atlas.Room, matched bool) {
	if !botutil.HasLinePrefix(s, prRoomExits) {
		return
	}

	room = atlas.NewEmptyRoom()
	matched = true

	// TODO: Cover all this in a lot of tests
	roomExitsStartPos := strings.Index(s, prRoomExits)
	roomExitsEndPos := roomExitsStartPos + strings.Index(s[roomExitsStartPos:], "\n") - 1
	roomExitsStr := strings.Trim(botutil.StripAnsi(s[roomExitsStartPos+len(prRoomExits):roomExitsEndPos]), "\r\n] ")
	if roomExitsStr == "None" {
		roomExitsStr = ""
	}
	if roomExitsStr != "" {
		// TODO: Test closed doors case
		roomExistsStrArr := strings.Split(roomExitsStr, " ")
		for _, roomExitStr := range roomExistsStrArr {
			roomExitStrCanonical := p.CanonizeExit(roomExitStr)
			dir, dirOk := atlas.NewDirection(roomExitStrCanonical)
			if !dirOk {
				p.logger.Debugf("Unknown direction %v in room, skipping", roomExitStrCanonical)
				continue
			}
			room.Exits[dir] = 0
		}
	}

	// TODO: Handle error if room exits prefix goes first
	roomNameStartPos := strings.LastIndex(s[:roomExitsStartPos-1], roomNameMarker)
	if roomNameStartPos > -1 {
		roomNameEndPos := roomNameStartPos + strings.Index(s[roomNameStartPos:], "\n") - 1
		// TODO: Handle error if no room name end
		room.Name = botutil.StripAnsi(s[roomNameStartPos:roomNameEndPos])
		room.Description = botutil.StripAnsi(strings.Trim(s[roomNameEndPos+2:roomExitsStartPos], "\r\n "))
	} else {
		room.PartialInfo = true
	}

	roomItemsStartPos := roomExitsEndPos + strings.Index(s[roomExitsEndPos:], roomItemsMarker)
	roomMobsStartPos := roomItemsStartPos + strings.Index(s[roomItemsStartPos:], roomMobsMarker)
	roomItemsStr := strings.Trim(botutil.StripAnsi(s[roomItemsStartPos:roomMobsStartPos]), "\r\n")
	if roomItemsStr != "" {
		room.Items = strings.Split(roomItemsStr, "\r\n")
	} else {
		room.Items = []string{}
	}

	roomMobsEndPos := roomMobsStartPos + strings.Index(s[roomMobsStartPos:], roomMobsEndMarker)
	roomMobsStr := strings.Trim(botutil.StripAnsi(s[roomMobsStartPos:roomMobsEndPos]), "\r\n")
	if roomMobsStr != "" {
		room.Mobs = strings.Split(roomMobsStr, "\r\n")
	} else {
		room.Mobs = []string{}
	}

	return
}

func (p *Parser) CanonizeExit(exit string) (res string) {
	res = exit
	res = strings.ReplaceAll(res, "(", "")
	res = strings.ReplaceAll(res, ")", "")
	return
}
