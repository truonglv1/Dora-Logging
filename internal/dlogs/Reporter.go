package dlogs

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/Dora-Logging/internal/djson"
	. "github.com/Dora-Logging/internal/utils"
	"github.com/Dora-Logging/utils"
	"github.com/marpaia/graphite-golang"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	ReportMetric     = `stats.gauges.%v.dora.log.report.%v`
	ReportMetricUser = `stats.gauges.%v.dora.log.report.user.%v`
)

func (dl *DLog) reportLogging(hostname string) {

	for {
		metrics := make([]graphite.Metric, 0)

		totalOldUser := len(dl.loadAllOldUser())
		totalNewUser := len(dl.loadAllNewUser())
		totalUser := totalOldUser + totalNewUser

		totalActivedUserInToday := len(dl.loadAllActivedUserInRangeDay(0))
		totalActivedUserIn7Day := len(dl.loadAllActivedUserInRangeDay(7))
		totalActivedUserIn15Day := len(dl.loadAllActivedUserInRangeDay(15))

		nameTotalUser := fmt.Sprintf(ReportMetricUser, hostname, "num-total-user")
		metrics = append(metrics, graphite.NewMetric(nameTotalUser, strconv.Itoa(totalUser), time.Now().Unix()))

		nameTotalNewUser := fmt.Sprintf(ReportMetricUser, hostname, "num-new-user")
		metrics = append(metrics, graphite.NewMetric(nameTotalNewUser, strconv.Itoa(totalNewUser), time.Now().Unix()))

		nameTotalActivedUserInToday := fmt.Sprintf(ReportMetricUser, hostname, "num-actived-user-today")
		metrics = append(metrics, graphite.NewMetric(nameTotalActivedUserInToday, strconv.Itoa(totalActivedUserInToday), time.Now().Unix()))

		nameTotalActivedUserIn7Day := fmt.Sprintf(ReportMetricUser, hostname, "num-actived-user-in-7-day")
		metrics = append(metrics, graphite.NewMetric(nameTotalActivedUserIn7Day, strconv.Itoa(totalActivedUserIn7Day), time.Now().Unix()))

		nameTotalActivedUserIn15day := fmt.Sprintf(ReportMetricUser, hostname, "num-actived-user-in-15-day")
		metrics = append(metrics, graphite.NewMetric(nameTotalActivedUserIn15day, strconv.Itoa(totalActivedUserIn15Day), time.Now().Unix()))

		// map: [user: so luot action]
		num_actived_user_ios := make(map[string]int64)
		num_actived_user_android := make(map[string]int64)

		map_event_follow_user := make(map[int]map[string]bool)
		map_event_follow_user[2001] = make(map[string]bool)
		map_event_follow_user[2002] = make(map[string]bool)

		totalActionIos := 0
		totalActionAndroid := 0

		num_read_summary := 0
		num_read_detail := 0
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
				num_read_summary += 1
			} else if v.EventApp == 2002 {
				num_read_detail += 1
			}

			if v.OsGroup.OsCode == 7 {
				totalActionIos += 1
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

		num_total_action_summary_detail := num_read_summary + num_read_detail

		num_user_read_summary := len(map_event_follow_user[2001])
		num_user_read_detail := len(map_event_follow_user[2002])
		total_user_read_summary_or_detail := num_user_read_summary + num_user_read_detail

		nameUserAndroid := fmt.Sprintf(ReportMetric, hostname, "num-actived-user-android")
		metrics = append(metrics, graphite.NewMetric(nameUserAndroid, strconv.Itoa(totalUserAndroid), time.Now().Unix()))

		nameUserIos := fmt.Sprintf(ReportMetric, hostname, "num-actived-user-ios")
		metrics = append(metrics, graphite.NewMetric(nameUserIos, strconv.Itoa(totalUserIos), time.Now().Unix()))

		nameActionAndroid := fmt.Sprintf(ReportMetric, hostname, "num-action-android")
		metrics = append(metrics, graphite.NewMetric(nameActionAndroid, strconv.Itoa(totalActionAndroid), time.Now().Unix()))

		nameActionIos := fmt.Sprintf(ReportMetric, hostname, "num-action-ios")
		metrics = append(metrics, graphite.NewMetric(nameActionIos, strconv.Itoa(totalActionIos), time.Now().Unix()))

		// user read summary / total user read sum + detail
		ratioUserReadSummaryMetric := fmt.Sprintf(ReportMetric, hostname, "ratio_user_read_summary")
		ratioUserReadSummary := fmt.Sprintf("%f", math.Round(float64(num_user_read_summary)/float64(total_user_read_summary_or_detail)*100))
		metrics = append(metrics, graphite.NewMetric(ratioUserReadSummaryMetric, ratioUserReadSummary, time.Now().Unix()))

		ratioUserReadDetailMetric := fmt.Sprintf(ReportMetric, hostname, "ratio_user_read_detail")
		ratioUserReadDetail := fmt.Sprintf("%f", math.Round(float64(num_user_read_detail)/float64(total_user_read_summary_or_detail)*100))
		metrics = append(metrics, graphite.NewMetric(ratioUserReadDetailMetric, ratioUserReadDetail, time.Now().Unix()))

		// action read summary / total action
		ratioActionReadSummaryMetric := fmt.Sprintf(ReportMetric, hostname, "ratio_action_read_summary")
		ratioActionReadSummary := fmt.Sprintf("%f", math.Round(float64(num_read_summary)/float64(num_total_action_summary_detail)*100))
		metrics = append(metrics, graphite.NewMetric(ratioActionReadSummaryMetric, ratioActionReadSummary, time.Now().Unix()))

		ratioActionReadDetailMetric := fmt.Sprintf(ReportMetric, hostname, "ratio_action_read_detail")
		ratioActionReadDetail := fmt.Sprintf("%f", float64(num_read_detail)/float64(num_total_action_summary_detail)*100)
		metrics = append(metrics, graphite.NewMetric(ratioActionReadDetailMetric, ratioActionReadDetail, time.Now().Unix()))

		fmt.Println("report info users ================")
		fmt.Println("total user in sys: ", totalUser)
		fmt.Println("num new user in today: ", totalNewUser)
		fmt.Println("num actived user in today: ", totalActivedUserInToday)
		fmt.Println("num actived user in 7 day: ", totalActivedUserIn7Day)
		fmt.Println("num actived user in 15 day: ", totalActivedUserIn15Day)

		fmt.Println("report log in today ==============")
		fmt.Println("num user ios: ", totalUserIos)
		fmt.Println("num action user ios: ", totalActionIos)
		fmt.Println("num user android: ", totalUserAndroid)
		fmt.Println("num action user android: ", totalActionAndroid)
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

func (dl *DLog) backUpManagerUsers() {
	session, err := mgo.Dial("topica.ai:27017")
	session.DB("dora").Login("sontc", "congson@123")

	if err != nil {
		HandleError(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	context := session.DB("dora").C("users_logging")

	//read folder logging
	files, err := ioutil.ReadDir("log-back-up")
	if err != nil {
		log.Fatal(err)
	}

	users := make(map[string]djson.UsersLog)

	for _, file := range files {
		pathFile := fmt.Sprintf("%v%v", "log-back-up/", file.Name())
		f, err := os.Open(pathFile)
		if err != nil {
			utils.HandleError(err)
		}
		s := bufio.NewScanner(f)
		for s.Scan() {
			var v djson.ActionLog
			if err := json.Unmarshal(s.Bytes(), &v); err != nil {
				utils.HandleError(err)
			}
			userID := strings.Split(v.SessionId, "_")[0]
			if _, ok := users[userID]; !ok && len(userID) > 0 {
				users[userID] = djson.UsersLog{UserId: userID, TimeCreate: v.TimeCreate}
			}
		}
	}
	//insert
	for _, val := range users {
		insert := make(map[string]interface{})
		insert["_id"] = val.UserId
		insert["user_id"] = val.UserId
		insert["time_create"] = val.TimeCreate
		err = context.Insert(insert)
		if err != nil {
			HandleError(err)
		}
	}

}

//compare with today
func (dl *DLog) loadAllOldUser() []map[string]interface{} {
	session, err := mgo.Dial("topica.ai:27017")
	session.DB("dora").Login("sontc", "congson@123")

	if err != nil {
		HandleError(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	context := session.DB("dora").C("users")
	var results []map[string]interface{}

	timeBeginDay := utils.GetTimeBeginDay(time.Now())
	err = context.Find(bson.M{
		"created_time": bson.M{
			"$lt": timeBeginDay,
		},
	}).Select(bson.M{"_id": 1}).All(&results)
	if err != nil {
		HandleError(err)
	}
	return results
}

//compare with today
func (dl *DLog) loadAllNewUser() []map[string]interface{} {
	session, err := mgo.Dial("topica.ai:27017")
	session.DB("dora").Login("sontc", "congson@123")

	if err != nil {
		HandleError(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	context := session.DB("dora").C("users")
	var results []map[string]interface{}

	timeBeginDay := utils.GetTimeBeginDay(time.Now())
	err = context.Find(bson.M{
		"created_time": bson.M{
			"$gt": timeBeginDay,
		},
	}).Select(bson.M{"_id": 1}).All(&results)
	if err != nil {
		HandleError(err)
	}
	return results
}

//get users actived in a range day
func (dl *DLog) loadAllActivedUserInRangeDay(numDay int) []map[string]interface{} {
	session, err := mgo.Dial("topica.ai:27017")
	session.DB("dora").Login("sontc", "congson@123")

	if err != nil {
		HandleError(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	context := session.DB("dora").C("users")
	var results []map[string]interface{}

	timeBeginAfterNumDay := utils.GetTimeBeginRangeDay(numDay)
	err = context.Find(bson.M{
		"updated_time": bson.M{
			"$gt": timeBeginAfterNumDay,
		},
	}).Select(bson.M{"_id": -1}).All(&results)
	if err != nil {
		HandleError(err)
	}
	return results
}

func (dl *DLog) insertNewUser(users map[string]djson.UsersLog) {
	session, err := mgo.Dial("topica.ai:27017")
	session.DB("dora").Login("sontc", "congson@123")

	if err != nil {
		HandleError(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	context := session.DB("dora").C("users_logging")

	//insert
	for _, val := range users {
		insert := make(map[string]interface{})
		insert["_id"] = val.UserId
		insert["user_id"] = val.UserId
		insert["time_create"] = val.TimeCreate
		err = context.Insert(insert)
		if err != nil {
			HandleError(err)
		}
	}
}
