package server

import (
	_ "embed"
	"mudbot/botutil"
	"mudbot/static"
	"os"
)

// Relative to binary workdir
var devTplFilename = "./static/atlas_tpl.html"
var devTplDoNotRead bool

func (s *Server) getDevTpl() (devTpl string, err error) {
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

func (s *Server) getHtmlTemplate() string {
	tpl, devTplErr := s.getDevTpl()
	if devTplErr != nil {
		s.logger.Errorf("Cannot check dev tpl existence: %v", devTplErr)
	} else if tpl == "" {
		tpl = static.AtlasHtmlTpl
	}

	return tpl
}
