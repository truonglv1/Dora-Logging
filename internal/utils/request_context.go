package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/kinghub-gateway/internal/client"
)

const (
	sizeBuffer int = 1024
)

func GetBodyToString(c *gin.Context) string {
	body := make([]byte, 0)
	buf := make([]byte, sizeBuffer)
	for {
		num, _ := c.Request.Body.Read(buf)
		if num == 0 {
			break
		}
		body = append(body, buf[0:num]...)
		if num < sizeBuffer {
			break
		}
	}
	return string(body)
}

func GetClientInfo(c *gin.Context) client.ClientInfo {
	info := client.ClientInfo{}
	info.Ip = c.ClientIP()
	return info
}
