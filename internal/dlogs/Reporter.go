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

		map_event_follow_user := make(map[int]map[string]bool)
		map_event_follow_user[2001] = make(map[string]bool)
		map_event_follow_user[2002] = make(map[string]bool)

		totalActionIOs := 0
		totalActionAndroid := 0

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
			if _, ok := map_event_follow_user[v.EventApp]; ok {
				map_event_follow_user[v.EventApp][userid] = true
			}

			if v.EventApp == 2001 {
				num_readed_summary += 1
			} else if v.EventApp == 2002 {
				num_readed_detail += 1
			}

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

		num_actived_total_user := len(num_actived_user_android) + len(num_actived_user_ios)
		num_total_action_summary_detail := num_readed_summary + num_readed_detail

		num_user_readed_summary := len(map_event_follow_user[2001])
		num_user_readed_detail := len(map_event_follow_user[2002])

		metrics := make([]graphite.Metric, 0)

		nameUserAndroid := fmt.Sprintf(ReportMetric, hostname, "num-actived-user-android")
		metrics = append(metrics, graphite.NewMetric(nameUserAndroid, strconv.Itoa(totalUserAndroid), time.Now().Unix()))

		nameUserIos := fmt.Sprintf(ReportMetric, hostname, "num-actived-user-ios")
		metrics = append(metrics, graphite.NewMetric(nameUserIos, strconv.Itoa(totalUserIos), time.Now().Unix()))

		nameActionAndroid := fmt.Sprintf(ReportMetric, hostname, "num-action-android")
		metrics = append(metrics, graphite.NewMetric(nameActionAndroid, strconv.Itoa(totalActionAndroid), time.Now().Unix()))

		nameActionIos := fmt.Sprintf(ReportMetric, hostname, "num-action-ios")
		metrics = append(metrics, graphite.NewMetric(nameActionIos, strconv.Itoa(totalActionIOs), time.Now().Unix()))

		ratioUserReadSummaryMetric := fmt.Sprintf(ReportMetric, hostname, "ratio_user_read_summary") // user read summary / total user
		ratioUserReadSummary := int(float64(num_user_readed_summary) / float64(num_actived_total_user) * 100)
		metrics = append(metrics, graphite.NewMetric(ratioUserReadSummaryMetric, strconv.Itoa(ratioUserReadSummary), time.Now().Unix()))

		ratioUserReadDetailMetric := fmt.Sprintf(ReportMetric, hostname, "ratio_user_read_detail")
		ratioUserReadDetail := int(float64(num_user_readed_detail) / float64(num_actived_total_user) * 100)
		metrics = append(metrics, graphite.NewMetric(ratioUserReadDetailMetric, strconv.Itoa(ratioUserReadDetail), time.Now().Unix()))

		ratioActionReadSummaryMetric := fmt.Sprintf(ReportMetric, hostname, "ratio_action_read_summary") // action read summary / total action
		ratioActionReadSummary := int(float64(num_readed_summary) / float64(num_total_action_summary_detail) * 100)
		metrics = append(metrics, graphite.NewMetric(ratioActionReadSummaryMetric, strconv.Itoa(ratioActionReadSummary), time.Now().Unix()))

		ratioActionReadDetailMetric := fmt.Sprintf(ReportMetric, hostname, "ratio_action_read_detail")
		ratioActionReadDetail := int(float64(num_readed_detail) / float64(num_total_action_summary_detail) * 100)
		metrics = append(metrics, graphite.NewMetric(ratioActionReadDetailMetric, strconv.Itoa(ratioActionReadDetail), time.Now().Unix()))

		fmt.Println("total user ios: ", totalUserIos)
		fmt.Println("total action user ios: ", totalActionIOs)
		fmt.Println("total user android: ", totalUserAndroid)
		fmt.Println("total action user android: ", totalActionAndroid)
		fmt.Println("ratio user read summary: ", ratioUserReadSummary)
		fmt.Println("ratio user read detail: ", ratioUserReadDetail)
		fmt.Println("ratio action read summary: ", ratioActionReadSummary)
		fmt.Println("ratio action read detail: ", ratioActionReadDetail)

		fmt.Println("============================")

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
