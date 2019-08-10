package dlogs

import (
	"github.com/gin-gonic/gin"
)

type DLog struct {
	router  *gin.Engine
	conf    *Config
	logChan chan Tuple
}

type Config struct {
	ServerAddr string
	ModeDebug  int
}

type Tuple struct {
	path string
	data map[string][]string
}
