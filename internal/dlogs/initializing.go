package dlogs

import (
	"fmt"
	fc "github.com/dora-logs/utils"
	"github.com/gin-gonic/gin"
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
	dl.router = gin.Default()
	dl.router.GET("/", dl.home)

	//api
	apiLog := dl.router.Group("/logging")
	_ = apiLog
	apiLog.GET("/trace", dl.trace)
	apiLog.POST("/trace", dl.tracePost)
	apiLog.POST("/trace/dev", dl.tracePostNew)

}
