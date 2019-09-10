package dlogs

import (
	"github.com/Dora-Logs/internal/djson"
	"github.com/gin-gonic/gin"
	"github.com/marpaia/graphite-golang"
)

type DLog struct {
	router  *gin.Engine
	conf    *Config
	logChan chan Tuple

	//metrics
	graphite *graphite.Graphite
	//counterAspect *metrics.CounterAspect
}

type Config struct {
	ServerAddr string
	ModeDebug  int
}

type Tuple struct {
	path      string
	actionLog []djson.ActionLog
}
