package tests

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/Dora-Logging/internal/djson"
	"github.com/Dora-Logging/utils"
	logj4 "github.com/jeanphorn/log4go"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
	"time"
)

func TestSaveLog(t *testing.T) {
	logj4.LoadConfiguration("./example.json")
	logj4.LOGGER("TestRotate").Info("category Test info test ...")
	logj4.LOGGER("Test").Info("category Test info test message: %s", "new test msg")
	logj4.LOGGER("Test").Debug("category Test debug test ...")

	// Other category not exist, test
	logj4.LOGGER("Other").Debug("category Other debug test ...")

	// socket log test
	logj4.LOGGER("TestSocket").Debug("category TestSocket debug test ...")

	// original log4go test
	logj4.Info("normal info test ...")
	logj4.Debug("normal debug test ...")

	logj4.Close()
}

func TestReadFile(t *testing.T) {
	f, err := os.Open("../log.log")
	if err != nil {
		utils.HandleError(err)
	}
	s := bufio.NewScanner(f)
	for s.Scan() {
		var v djson.ActionLog
		if err := json.Unmarshal(s.Bytes(), &v); err != nil {
			utils.HandleError(err)
		}
		fmt.Println(v.Ip)
	}
	if s.Err() != nil {
		utils.HandleError(err)
	}

	//file, _ := ioutil.ReadFile("log.log")
	//data := []djson.ActionLog{}
	//err := json.Unmarshal([]byte(file), &data)
	//if err != nil {
	//	HandleError(err)
	//} else {
	//	for _, val := range data {
	//		fmt.Println(val)
	//	}
	//}
}

func TestParseJson(t *testing.T) {
	s := `{"ip":"1.0.0.1","os_group":{"os_code":8,"os_ver":"9","user_agent":"samsungSM-G950N"},"session_id":"test","category_id":"abcxyz","event_app":10000,"event_id":"","article_id":0,"time_create":11111}`
	data := &djson.ActionLog{}
	err := json.Unmarshal([]byte(s), data)
	fmt.Println(err)
	s2, _ := json.Marshal(data)
	fmt.Println(string(s2))
	fmt.Println(data.Ip)
}

func TestReadFolder(t *testing.T) {
	files, err := ioutil.ReadDir("../logging")
	if err != nil {
		log.Fatal(err)
	}

	for index, file := range files {
		fmt.Println(index)
		fmt.Println(file.Name())
		a := fmt.Sprintf("%v%v", "log-back-up/", file.Name())
		fmt.Println(a)
	}
}

func TestTime(t *testing.T) {
	now := time.Now()

	fmt.Println("now:", now)

	then := now.AddDate(0, 0, -1)

	fmt.Println("then:", then.Unix())
}

func TestLogg(t *testing.T) {
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

	f, err := os.Open("../logging/log.log")
	if err != nil {
		utils.HandleError(err)
	}
	s := bufio.NewScanner(f)

	for s.Scan() {
		var v djson.ActionLog
		if err := json.Unmarshal(s.Bytes(), &v); err != nil {
			utils.HandleError(err)
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

	num_user_readed_summary := len(map_event_follow_user[2001])
	num_user_readed_detail := len(map_event_follow_user[2002])

	ratioUserReadSummary := int(float64(num_user_readed_summary) / float64(num_actived_total_user) * 100)
	fmt.Println(num_user_readed_summary)
	fmt.Println(num_actived_total_user)
	fmt.Println(ratioUserReadSummary)
	fmt.Println("=============================")
	ratioUserReadDetail := int(float64(num_user_readed_detail) / float64(num_actived_total_user) * 100)
	fmt.Println(num_user_readed_detail)
	fmt.Println(num_actived_total_user)
	fmt.Println(ratioUserReadDetail)

	//ratioActionReadDetail := float64(num_readed_detail/num_total_action) * 100

	fmt.Println("=============================")
	fmt.Println("total user ios: ", totalUserIos)
	fmt.Println("total action user ios: ", totalActionIOs)
	fmt.Println("total user android: ", totalUserAndroid)
	fmt.Println("total action user android: ", totalActionAndroid)
	fmt.Println("ratio user read summary: ", ratioUserReadDetail)
	fmt.Println("============================")

}

func TestA(t *testing.T) {
	fmt.Println(float64(3) / float64(10))
}
