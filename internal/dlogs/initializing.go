package dlogs

import (
	"fmt"
	"github.com/Dora-Logging/internal/metrics"
	"github.com/Dora-Logging/internal/utils"
	fc "github.com/Dora-Logging/utils"
	"github.com/gin-gonic/gin"
	"github.com/marpaia/graphite-golang"
	"gopkg.in/mgo.v2"
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

	//load categories
	errCate := dl.initCategories()
	if errCate != nil{
		return nil, errCate
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

func (dl *DLog) initCategories() error   {
	dl.Categories = make(map[string]string)

	session, err := mgo.Dial("110.35.75.40:27072")
	session.DB("dora").Login("sontc", "c0ngSon@@3921")

	if err != nil {
		utils.HandleError(err)
		return err
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	context := session.DB("dora").C("categories")

	results := []Category{}
	e := context.Find(nil).All(&results)
	if e != nil {
		// TODO: Do something about the error
		utils.HandleError(e)
		return e
	} else {
		for _,c := range results  {
			dl.Categories[c.CategoryId.Hex()] = c.CategorySlug
		}
	}

	return nil
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
	logFile.SetOutput(outputError)

	outputFile := &lumberjack.Logger{
		Filename:   "/home/sontc/truonglv/Dora-Logging/server-logs/" + hostname + "-server.log",
		MaxSize:    128, // megabytes
		MaxBackups: 2,
		MaxAge:     7, //days
	}
	dl.router.Use(gin.LoggerWithWriter(outputFile))
	dl.router.Use(gin.Recovery())

	//initialize CounterAspect and reset every minute
	dl.graphite, _ = graphite.NewGraphite("110.35.75.40", 2003)

	////counter
	dl.counterAspect = metrics.NewCounterAspect(dl.graphite, hostname, dl.Categories)
	dl.counterAspect.StartTimer(10 * time.Minute)
	dl.router.Use(metrics.CounterHandler(dl.counterAspect))

	dl.router.GET("/", dl.home)

	//api logging app
	apiLogApp := dl.router.Group("/logging")
	//_ = apiLogApp
	apiLogApp.POST("/trace", dl.tracePost)
	apiLogApp.POST("/trace/dev", dl.tracePostNew)

	//go dl.reportLogging(hostname)
	//dl.loadAllActivedUserInRangeDay(7)

	//api logging web
	apiLogWeb := dl.router.Group("/web/logging")
	apiLogWeb.POST("/trace", dl.loggingOnWeb)

	apiLogWeb.GET("/trace", dl.loggingOnWeb)

}
