package dlogs

import (
	"encoding/hex"
	"github.com/gin-gonic/gin"
)

func (dl *DLog) home(c *gin.Context) {
	winNoticeImg, _ := hex.DecodeString("47494638396101000100800000" +
		"FFFFFF0000002C000000000100010000" +
		"02024401003B")
	c.Header("Content-Type", "image/gif")
	_, _ = c.Writer.Write(winNoticeImg)
}

func (dl *DLog) trace(c *gin.Context) {

	params := c.Request.URL.Query()
	dl.saveLog("trace", params)
	winNoticeImg, _ := hex.DecodeString("47494638396101000100800000" +
		"FFFFFF0000002C000000000100010000" +
		"02024401003B")
	c.Header("Content-Type", "image/gif")
	_, _ = c.Writer.Write(winNoticeImg)

}

//func (dl *DLog) activeApp(c *gin.Context)  {
//	params := c.Request.URL.Query()
//	dl.saveLog("activate", params)
//	winNoticeImg, _ := hex.DecodeString("47494638396101000100800000" +
//		"FFFFFF0000002C000000000100010000" +
//		"02024401003B")
//	c.Header("Content-Type", "image/gif")
//	_, _ = c.Writer.Write(winNoticeImg)
//
//}
