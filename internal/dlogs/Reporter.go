package dlogs

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/Dora-Logging/internal/djson"
	. "github.com/Dora-Logging/internal/utils"
	"github.com/marpaia/graphite-golang"
	"os"
	"strconv"
	"time"
)

const (
	ReportMetric = `stats.gauges.%v.dora.log.report.%v`
)

func (dl *DLog) reportLogging(hostname string) {

	for {
		// map: [user: so luot action]
		listUserIOS := make(map[string]int64)
		listUserAndroid := make(map[string]int64)
		totalActionIOs := 0
		totalActionAndroid := 0

		f, err := os.Open("logging/log.log")
		if err != nil {
			HandleError(err)
		}
		s := bufio.NewScanner(f)
		for s.Scan() {
			var v djson.ActionLog
			if err := json.Unmarshal(s.Bytes(), &v); err != nil {
				HandleError(err)
			}
			if v.OsGroup.OsCode == 7 {
				totalActionIOs += 1
				if _, ok := listUserIOS[v.OsGroup.UserAgent]; !ok {
					listUserIOS[v.OsGroup.UserAgent] = 1
				}
			} else if v.OsGroup.OsCode == 8 {
				totalActionAndroid += 1
				if _, ok := listUserAndroid[v.OsGroup.UserAgent]; !ok {
					listUserAndroid[v.OsGroup.UserAgent] = 1
				}
			}
		}
		totalUserIos := len(listUserIOS)
		totalUserAndroid := len(listUserAndroid)

		fmt.Println("total user ios: ", totalUserIos)
		fmt.Println("total action user ios: ", totalActionIOs)
		fmt.Println("total user android: ", totalUserAndroid)
		fmt.Println("total action user android: ", totalActionAndroid)
		fmt.Println("============================")

		metrics := make([]graphite.Metric, 0)

		nameUserAndroid := fmt.Sprintf(ReportMetric, hostname, "total-user-android")
		metrics = append(metrics, graphite.NewMetric(nameUserAndroid, strconv.Itoa(totalUserAndroid), time.Now().Unix()))

		nameUserIos := fmt.Sprintf(ReportMetric, hostname, "total-user-ios")
		metrics = append(metrics, graphite.NewMetric(nameUserIos, strconv.Itoa(totalUserIos), time.Now().Unix()))

		nameActionAndroid := fmt.Sprintf(ReportMetric, hostname, "total-action-android")
		metrics = append(metrics, graphite.NewMetric(nameActionAndroid, strconv.Itoa(totalActionAndroid), time.Now().Unix()))

		nameActionIos := fmt.Sprintf(ReportMetric, hostname, "total-action-ios")
		metrics = append(metrics, graphite.NewMetric(nameActionIos, strconv.Itoa(totalActionIOs), time.Now().Unix()))
		//send metrics
		if len(metrics) > 0 {
			err := dl.graphite.SendMetrics(metrics)
			if err != nil {
				HandleError(err)
			}
		}
		if s.Err() != nil {
			HandleError(err)
		}
		time.Sleep(1 * time.Hour)
		fmt.Println(time.Now().Date())
	}
}
