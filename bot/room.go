package bot

import (
	"mudbot/botutil"
	"strings"
)

var prMapper = "\x1b[0;36m[ Exits: "
var roomNameMarker = "\x1b[1;36m"
var roomItemsMarker = "\x1b[1;33m"
var roomMobsMarker = "\x1b[1;31m"
var roomMobsEndMarker = "\x1b[0;0m"

func (b *Bot) ParseRoom(s string) Event {
	if !botutil.HasLinePrefix(s, prMapper) {
		return EVENT_NOP
	}

	// TODO: Cover all this in a lot of tests
	mapperStartPos := strings.Index(s, prMapper)
	mapperEndPos := mapperStartPos + strings.Index(s[mapperStartPos:], "\n")
	// TODO: Handle error if mapper prefix goes first
	roomNameStartPos := strings.LastIndex(s[:mapperStartPos-1], roomNameMarker)
	if roomNameStartPos > -1 {
		roomNameEndPos := roomNameStartPos + strings.Index(s[roomNameStartPos:], "\n")
		// TODO: Handle error if no room name end
		roomName := botutil.StripAnsi(s[roomNameStartPos:roomNameEndPos])
		b.logger.Debugf("room name: %v", roomName)
		b.logger.Debugf("room name bytes:\n%v", botutil.StrToHex(roomName))
		roomDesc := botutil.StripAnsi(strings.Trim(s[roomNameEndPos+2:mapperStartPos], "\n "))
		b.logger.Debugf("room desc: %v", roomDesc)
	}
	roomItemsStartPos := mapperEndPos + strings.Index(s[mapperEndPos:], roomItemsMarker)
	roomMobsStartPos := roomItemsStartPos + strings.Index(s[roomItemsStartPos:], roomMobsMarker)
	roomItemsStr := botutil.StripAnsi(strings.Trim(s[roomItemsStartPos:roomMobsStartPos], "\n"))
	roomItems := strings.Split(roomItemsStr, "\n")
	b.logger.Debugf("room items (%n):\n%+v", len(roomItems), roomItems)
	b.logger.Debugf("room items bytes:\n%v", roomItemsStr)

	roomMobsEndPos := roomMobsStartPos + strings.Index(s[roomMobsStartPos:], roomMobsEndMarker)
	roomMobsStr := botutil.StripAnsi(strings.Trim(s[roomMobsStartPos:roomMobsEndPos], "\n"))
	roomMobs := strings.Split(roomMobsStr, "\n")
	b.logger.Debugf("room mobs (%n):\n%+v", len(roomMobs), roomMobs)
	b.logger.Debugf("room mobs bytes:\n%v", roomMobsStr)

	return EVENT_ROOM
}
