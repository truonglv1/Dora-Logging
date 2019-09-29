package dlogs

import (
	"github.com/Dora-Logs/internal/djson"
	"github.com/Dora-Logs/internal/metrics"
	"github.com/gin-gonic/gin"
	"github.com/marpaia/graphite-golang"
)

type DLog struct {
	router  *gin.Engine
	conf    *Config
	logChan chan Tuple

	//metrics
	graphite      *graphite.Graphite
	counterAspect *metrics.CounterAspect

	report *Report
}

type Config struct {
	ServerAddr string
	ModeDebug  int
}

type Tuple struct {
	path      string
	actionLog []djson.ActionLog
}

type Report struct {
	numberUserIos      int64 `json:"number_user_ios"`
	numberUserAndroid  int64 `json:"number_user_android"`
	totalActionIos     int64 `json:"total_action_ios"`
	totalActionAndroid int64 `json:"total_action_android"`
}
