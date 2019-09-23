package dlogs

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/Dora-Logging/internal/djson"
	. "github.com/Dora-Logging/internal/utils"
	"github.com/marpaia/graphite-golang"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	ReportMetric = `stats.gauges.%v.dora.log.report.%v`
)

func (dl *DLog) reportLogging(hostname string) {

	for {
		// map: [user: so luot action]
		num_actived_user_ios := make(map[string]int64)
		num_actived_user_android := make(map[string]int64)

		totalActionIOs := 0
		totalActionAndroid := 0

		num_user_readed_summary := 0
		num_user_readed_detail := 0

		num_readed_summary := 0
		num_readed_detail := 0
		//ti le user doc bai tóm bắt (event: 2001), bài chi tiết (event: 2002)

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
			userid := strings.Split(v.SessionId, "_")[0]

			if v.EventApp == 2001 {
				num_readed_summary += 1
			} else if v.EventApp == 2002 {
				num_readed_detail += 1
			}

			if v.OsGroup.OsCode == 7 {
				totalActionIOs += 1
				if _, ok := num_actived_user_ios[userid]; !ok {
					num_actived_user_ios[userid] = 1
					if v.EventApp == 2001 {
						num_user_readed_summary += 1
					} else if v.EventApp == 2002 {
						num_user_readed_detail += 1
					}
				}
			} else if v.OsGroup.OsCode == 8 {
				totalActionAndroid += 1
				if _, ok := num_actived_user_android[userid]; !ok {
					num_actived_user_android[userid] = 1
					if v.EventApp == 2001 {
						num_user_readed_summary += 1
					} else if v.EventApp == 2002 {
						num_user_readed_detail += 1
					}
				}
			}

		}
		totalUserIos := len(num_actived_user_ios)
		totalUserAndroid := len(num_actived_user_android)

		fmt.Println("total user ios: ", totalUserIos)
		fmt.Println("total action user ios: ", totalActionIOs)
		fmt.Println("total user android: ", totalUserAndroid)
		fmt.Println("total action user android: ", totalActionAndroid)
		fmt.Println("============================")

		num_actived_total_user := len(num_actived_user_android) + len(num_actived_user_ios)
		num_total_action := totalActionAndroid + totalActionIOs

		metrics := make([]graphite.Metric, 0)

		nameUserAndroid := fmt.Sprintf(ReportMetric, hostname, "num-actived-user-android")
		metrics = append(metrics, graphite.NewMetric(nameUserAndroid, strconv.Itoa(totalUserAndroid), time.Now().Unix()))

		nameUserIos := fmt.Sprintf(ReportMetric, hostname, "num-actived-user-ios")
		metrics = append(metrics, graphite.NewMetric(nameUserIos, strconv.Itoa(totalUserIos), time.Now().Unix()))

		nameActionAndroid := fmt.Sprintf(ReportMetric, hostname, "num-action-android")
		metrics = append(metrics, graphite.NewMetric(nameActionAndroid, strconv.Itoa(totalActionAndroid), time.Now().Unix()))

		nameActionIos := fmt.Sprintf(ReportMetric, hostname, "num-action-ios")
		metrics = append(metrics, graphite.NewMetric(nameActionIos, strconv.Itoa(totalActionIOs), time.Now().Unix()))

		ratioUserReadSummary := fmt.Sprintf(ReportMetric, hostname, "ratio_user_readed_summary") // user read summary / total user
		metrics = append(metrics, graphite.NewMetric(ratioUserReadSummary, strconv.Itoa(num_user_readed_summary/num_actived_total_user), time.Now().Unix()))
		ratioUserReadDetail := fmt.Sprintf(ReportMetric, hostname, "ratio_user_readed_detail")
		metrics = append(metrics, graphite.NewMetric(ratioUserReadDetail, strconv.Itoa(num_user_readed_detail/num_actived_total_user), time.Now().Unix()))

		ratioActionReadSummary := fmt.Sprintf(ReportMetric, hostname, "ratio_action_readed_summary") // action read summary / total action
		metrics = append(metrics, graphite.NewMetric(ratioActionReadSummary, strconv.Itoa(num_readed_summary/num_total_action), time.Now().Unix()))
		ratioActionReadDetail := fmt.Sprintf(ReportMetric, hostname, "ratio_action_readed_detail")
		metrics = append(metrics, graphite.NewMetric(ratioActionReadDetail, strconv.Itoa(num_readed_detail/num_total_action), time.Now().Unix()))

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
		time.Sleep(10 * time.Minute)
		fmt.Println(time.Now().Date())
	}
}

func (dl *DLog) reportLoggingBackup(hostname string) {
	files, err := ioutil.ReadDir("log-back-up")
	if err != nil {
		log.Fatal(err)
	}
	now := time.Now()
	count := 5
	for _, file := range files {
		then := now.AddDate(0, 0, -count)
		count = -1
		fmt.Println(file.Name())
		// map: [user: so luot action]
		num_actived_user_ios := make(map[string]int64)
		num_actived_user_android := make(map[string]int64)
		totalActionIOs := 0
		totalActionAndroid := 0

		f, err := os.Open(fmt.Sprintf("%v%v", "log-back-up/", file.Name()))
		if err != nil {
			HandleError(err)
		}
		s := bufio.NewScanner(f)
		for s.Scan() {
			var v djson.ActionLog
			if err := json.Unmarshal(s.Bytes(), &v); err != nil {
				HandleError(err)
			}
			userid := strings.Split(v.SessionId, "_")[0]
			if v.OsGroup.OsCode == 7 {
				totalActionIOs += 1
				if _, ok := num_actived_user_ios[userid]; !ok {
					num_actived_user_ios[userid] = 1
				}
			} else if v.OsGroup.OsCode == 8 {
				totalActionAndroid += 1
				if _, ok := num_actived_user_android[userid]; !ok {
					num_actived_user_android[userid] = 1
				}
			}
		}
		totalUserIos := len(num_actived_user_ios)
		totalUserAndroid := len(num_actived_user_android)

		fmt.Println("total user ios: ", totalUserIos)
		fmt.Println("total action user ios: ", totalActionIOs)
		fmt.Println("total user android: ", totalUserAndroid)
		fmt.Println("total action user android: ", totalActionAndroid)
		fmt.Println("============================")

		metrics := make([]graphite.Metric, 0)

		nameUserAndroid := fmt.Sprintf(ReportMetric, hostname, "num-actived-user-android")
		metrics = append(metrics, graphite.NewMetric(nameUserAndroid, strconv.Itoa(totalUserAndroid), then.Unix()))

		nameUserIos := fmt.Sprintf(ReportMetric, hostname, "num-actived-user-ios")
		metrics = append(metrics, graphite.NewMetric(nameUserIos, strconv.Itoa(totalUserIos), then.Unix()))

		nameActionAndroid := fmt.Sprintf(ReportMetric, hostname, "num-action-android")
		metrics = append(metrics, graphite.NewMetric(nameActionAndroid, strconv.Itoa(totalActionAndroid), then.Unix()))

		nameActionIos := fmt.Sprintf(ReportMetric, hostname, "num-action-ios")
		metrics = append(metrics, graphite.NewMetric(nameActionIos, strconv.Itoa(totalActionIOs), then.Unix()))

		//send metrics
		if len(metrics) > 0 {
			fmt.Println("is sendingggg ....")
			err := dl.graphite.SendMetrics(metrics)
			if err != nil {
				HandleError(err)
			}
		}
		if s.Err() != nil {
			HandleError(err)
		}
	}

}
