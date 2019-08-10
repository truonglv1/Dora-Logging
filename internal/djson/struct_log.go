package djson

type ActionLog struct {
	Ip string `json:"ip"`
	Os string `json:"os"`
	SessionId   string `json:"session_id"`
	CategoryId string `json:"category_id"`
	EventId    int   `json:"event_id"`
	TimeCreate int64 `json:"time_create"` //utc0
}

type OsGroup struct {
	OsCode int `json:"os_code"`
	OsVer string `json:"os_ver"`
}
