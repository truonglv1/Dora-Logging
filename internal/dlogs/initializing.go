package dlogs

import (
	"fmt"
	"github.com/Dora-Logging/internal/metrics"
	fc "github.com/Dora-Logging/utils"
	"github.com/gin-gonic/gin"
	"github.com/marpaia/graphite-golang"
	"gopkg.in/natefinch/lumberjack.v2"
	logFile "log"
	"os"
	"time"
)

func InitServerLogging(pConf string) (*DLog, error) {
	if len(pConf) == 0 {
		pConf = pathConf
	}
	dl := &DLog{}
	err := dl.initConf(pConf)
	if err != nil {
		return nil, err
	}
	dl.initRoute()
	dl.initLog()
	return dl, nil
}

func (dl *DLog) ListenAndServe() {
	if dl.conf.ModeDebug == 0 {
		fmt.Printf("Listening and serving HTTP on %s\n", dl.conf.ServerAddr)
	}
	err := dl.router.Run(dl.conf.ServerAddr)
	if err != nil {
		panic(err)
	}
}

func (dl *DLog) initConf(pConf string) error {
	dl.conf = &Config{}
	return fc.LoadConfig(pathConf, dl.conf)
}

func (dl *DLog) initRoute() {
	if dl.conf.ModeDebug == 0 {
		gin.SetMode(gin.ReleaseMode)
	}
	dl.router = gin.New()
	hostname, _ := os.Hostname()
	outputError := &lumberjack.Logger{
		Filename:   "/home/sontc/truonglv/Dora-Logging/server-logs/" + hostname + "-error.log",
		MaxSize:    128, // megabytes
		MaxBackups: 2,
		MaxAge:     7, //days
	}

	outputFile := &lumberjack.Logger{
		Filename:   "/home/sontc/truonglv/Dora-Logging/server-logs/" + hostname + "-server.log",
		MaxSize:    128, // megabytes
		MaxBackups: 2,
		MaxAge:     7, //days
	}
	logFile.SetOutput(outputError)
	dl.router.Use(gin.LoggerWithWriter(outputFile))
	dl.router.Use(gin.Recovery())

	//initialize CounterAspect and reset every minute
	dl.graphite, _ = graphite.NewGraphite("42.113.206.204", 2003)
	////counter
	dl.counterAspect = metrics.NewCounterAspect(dl.graphite, hostname)
	dl.counterAspect.StartTimer(1 * time.Minute)
	dl.router.Use(metrics.CounterHandler(dl.counterAspect))

	dl.router.GET("/", dl.home)

	//api
	apiLog := dl.router.Group("/logging")
	_ = apiLog
	apiLog.GET("/trace", dl.trace)
	apiLog.POST("/trace", dl.tracePost)
	apiLog.POST("/trace/dev", dl.tracePostNew)

	go dl.reportLogging(hostname)

}
