package dlogs

import (
	"encoding/json"
	"fmt"
	"github.com/Dora-Logs/internal/djson"
	. "github.com/Dora-Logs/internal/utils"
	"io/ioutil"
)

func (dl *DLog) reportLogging() {
	fmt.Println("aaaaaa")
	file, _ := ioutil.ReadFile("log.log")
	data := []djson.ActionLog{}
	err := json.Unmarshal([]byte(file), &data)
	if err != nil {
		HandleError(err)
	} else {
		for _, val := range data {
			fmt.Println(val)
		}
	}
}
