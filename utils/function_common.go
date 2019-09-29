package utils

import (
	"github.com/BurntSushi/toml"
	"log"
	"os"
	"runtime"
	"time"
)

// config is must point
func LoadConfig(fileName string, config interface{}) error {
	_, err := os.Stat(fileName)
	if err != nil {
		return err
	}
	if _, err := toml.DecodeFile(fileName, config); err != nil {
		return err
	}
	return err
}

func HandleError(err interface{}) {
	if err != nil {
		_, fn, line, _ := runtime.Caller(1)
		log.Printf("[E] %v %s:%d", err, fn, line)
	}
}

func GetTimeBeginDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 918273645, time.UTC)
}

func GetTimeBeginRangeDay(num_day int) time.Time {
	year, month, day := time.Now().AddDate(0, 0, -num_day).Date()
	return time.Date(year, month, day, 0, 0, 0, 918273645, time.UTC)
}

func ContainSlice(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
