package dlogs

import (
	"encoding/json"
	"fmt"
	"github.com/Dora-Logging/internal/djson"
	logj4 "github.com/jeanphorn/log4go"
	"os"
	"os/signal"
	"syscall"
)

func (dl *DLog) initLog() {
	logj4.LoadConfiguration(pathLogConf)
	dl.startLog()
	dl.startLogWeb()
}

func (dl *DLog) startLog() {
	dl.logChan = make(chan Tuple)
	var gracefulStop = make(chan os.Signal, 1)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	go func() {
		for {
			select {
			case tup := <-dl.logChan:
				dl.printLog(tup)
			case <-gracefulStop:
				logj4.Close()
				os.Exit(0)
			}
		}
	}()
}

func (dl *DLog) startLogWeb() {
	dl.logChanWeb = make(chan TupleWeb)
	var gracefulStop = make(chan os.Signal, 1)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	go func() {
		for {
			select {
			case tup := <-dl.logChanWeb:
				dl.printLogWeb(tup)
			case <-gracefulStop:
				logj4.Close()
				os.Exit(0)
			}
		}
	}()
}

func (dl *DLog) printLog(tup Tuple) {

	for _, val := range tup.actionLog {
		b, err := json.Marshal(val)
		if err != nil {
			fmt.Println(err)
			return
		}
		logj4.LOGGER("app-log").Info(string(b))
	}
}

func (dl *DLog) printLogWeb(tup TupleWeb) {

	b, err := json.Marshal(tup.webAction)
	if err != nil {
		fmt.Println(err)
		return
	}
	logj4.LOGGER("web-log").Info(string(b))
}

func (dl *DLog) saveLog(path string, actionLogs []djson.ActionLog) {
	dl.logChan <- Tuple{
		path:      path,
		actionLog: actionLogs,
	}
}

func (dl *DLog) saveLogWeb(path string, actionLog djson.WebAction) {
	dl.logChanWeb <- TupleWeb{
		path:      path,
		webAction: actionLog,
	}
}
