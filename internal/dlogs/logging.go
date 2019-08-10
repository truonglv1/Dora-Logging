package dlogs

import (
	"encoding/json"
	"fmt"
	"github.com/Dora-Logs/internal/djson"
	logj4 "github.com/jeanphorn/log4go"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func (dl *DLog) initLog() {
	logj4.LoadConfiguration(pathLogConf)
	dl.startLog()
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

func (dl *DLog) printLog(tup Tuple) {
	a := &djson.ActionLog{}
	os := &djson.OsGroup{}

	a.TimeCreate = time.Now().UTC().Unix()
	if val, ok := tup.data["category_id"]; ok {
		a.CategoryId = val[0]
	}
	if val, ok := tup.data["event_id"]; ok {
		a.EventId, _ = strconv.Atoi(val[0])
	}

	if val, ok := tup.data["ip"]; ok {
		a.Ip = val[0]
	}
	if val, ok := tup.data["session_id"]; ok {
		a.SessionId = val[0]

	}

	if val, ok := tup.data["os_code"]; ok {
		os.OsCode, _ = strconv.Atoi(val[0])
	}

	if val, ok := tup.data["os_ver"]; ok {
		os.OsVer = val[0]
	}
	a.OsGroup = *os

	b, err := json.Marshal(a)
	if err != nil {
		fmt.Println(err)
		return
	}

	logj4.LOGGER("app-log").Info(string(b))
}

func (dl *DLog) saveLog(path string, data map[string][]string) {
	dl.logChan <- Tuple{
		path: path,
		data: data,
	}
}
