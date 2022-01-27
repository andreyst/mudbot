package atlas

import (
	_ "embed"
	"mudbot/botutil"
	"mudbot/static"
	"os"
)

var devTplFilename = "./static/atlas_tpl.html"
var devTplDoNotRead bool

func getDevTpl() (devTpl string, err error) {
	var devTplExists bool
	devTplExists, err = botutil.Exists(devTplFilename)
	if err != nil {
		devTplDoNotRead = true
		return
	}

	if !devTplDoNotRead && devTplExists {
		var devTplBytes []byte
		devTplBytes, err = os.ReadFile(devTplFilename)
		if err != nil {
			devTplDoNotRead = true
			return
		}
		devTpl = string(devTplBytes)
	}

	return
}

func (a Atlas) getHtmlTemplate() string {
	tpl, devTplErr := getDevTpl()
	if devTplErr != nil {
		a.logger.Errorf("Cannot check dev tpl existence: %v", devTplErr)
	} else if tpl == "" {
		tpl = static.AtlasHtmlTpl
	}

	//res := strings.ReplaceAll(tpl, "{rooms}", string(roomsJson))
	//res = strings.ReplaceAll(res, "{currentCoordinates}", string(currentCoordinatesJson))

	return tpl
}
