package tests

import (
	logj4 "github.com/jeanphorn/log4go"
	"testing"
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
