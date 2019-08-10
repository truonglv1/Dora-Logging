package djson

type ActionLog struct {
	Ip         string  `json:"ip"`
	OsGroup    OsGroup `json:"os_group"`
	SessionId  string  `json:"session_id"`
	CategoryId string  `json:"category_id"`
	EventId    int     `json:"event_id"`
	TimeCreate int64   `json:"time_create"` //utc0
}

type OsGroup struct {
	OsCode    int    `json:"os_code"`
	OsVer     string `json:"os_ver"`
	UserAgent string `json:"user_agent"`
}
