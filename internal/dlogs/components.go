package dlogs

import (
	"github.com/Dora-Logging/internal/djson"
	"github.com/Dora-Logging/internal/metrics"
	"github.com/gin-gonic/gin"
	"github.com/marpaia/graphite-golang"
	"gopkg.in/mgo.v2/bson"
)

type DLog struct {
	router  *gin.Engine
	conf    *Config
	logChan chan Tuple
	logChanWeb chan TupleWeb

	//metrics
	graphite      *graphite.Graphite
	counterAspect *metrics.CounterAspect

	// category
	Categories	map[string]string

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

type TupleWeb struct {
	path      string
	webAction djson.WebAction
}

type Report struct {
	numberUserIos      int64 `json:"number_user_ios"`
	numberUserAndroid  int64 `json:"number_user_android"`
	totalActionIos     int64 `json:"total_action_ios"`
	totalActionAndroid int64 `json:"total_action_android"`
}

type Category struct {
	CategoryId bson.ObjectId `json:"category_id" bson:"_id"`
	CategoryName string `json:"category_name" bson:"name"`
	CategorySlug string `json:"slug" bson:"slug"`
}