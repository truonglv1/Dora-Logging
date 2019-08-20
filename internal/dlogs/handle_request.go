package dlogs

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/Dora-Logs/internal/djson"
	"github.com/gin-gonic/gin"
	"io/ioutil"
)

func (dl *DLog) home(c *gin.Context) {
	winNoticeImg, _ := hex.DecodeString("47494638396101000100800000" +
		"FFFFFF0000002C000000000100010000" +
		"02024401003B")
	c.Header("Content-Type", "image/gif")
	_, _ = c.Writer.Write(winNoticeImg)
}

func (dl *DLog) trace(c *gin.Context) {

	//params := c.Request.URL.Query()
	//dl.saveLog("trace", params)
	//winNoticeImg, _ := hex.DecodeString("47494638396101000100800000" +
	//	"FFFFFF0000002C000000000100010000" +
	//	"02024401003B")
	//c.Header("Content-Type", "image/gif")
	//_, _ = c.Writer.Write(winNoticeImg)

}

func (dl *DLog) tracePost(c *gin.Context) {

	body, err := ioutil.ReadAll(c.Request.Body)
	defer c.Request.Body.Close()
	if err != nil {
		fmt.Println("err1: ", err)
		c.Writer.Write([]byte(err.Error()))
		return
	}
	actionLogs := []djson.ActionLog{}
	err = json.Unmarshal(body, &actionLogs)
	if err != nil {
		fmt.Println("err2: ", err)
		c.Writer.Write([]byte(err.Error()))
		return
	} else if len(actionLogs) > 0 {
		dl.saveLog("trace", actionLogs)
		winNoticeImg, _ := hex.DecodeString("47494638396101000100800000" +
			"FFFFFF0000002C000000000100010000" +
			"02024401003B")
		c.Header("Content-Type", "image/gif")
		_, _ = c.Writer.Write(winNoticeImg)
		return
	}
	c.Writer.Write([]byte("error"))
	return
}
