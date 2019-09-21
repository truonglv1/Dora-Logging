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
