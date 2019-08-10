package utils

import (
	"github.com/BurntSushi/toml"
	"log"
	"os"
	"runtime"
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
