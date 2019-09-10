package djson

type ActionLog struct {
	Ip         string  `json:"ip"`
	OsGroup    OsGroup `json:"os_group"`
	SessionId  string  `json:"session_id"`
	CategoryId string  `json:"category_id"`
	EventApp   int     `json:"event_app"`
	EventId    string  `json:"event_id"`
	ArticleId  int     `json:"article_id"`

	TimeCreate int64 `json:"time_create"` //utc0
}

type OsGroup struct {
	OsCode    int    `json:"os_code"`
	OsVer     string `json:"os_ver"`
	UserAgent string `json:"user_agent"`
}

type ResponseClient struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Code    interface{} `json:"code"`
	Data    interface{} `json:"data,omitempty"`
}

type Data_res struct {
	SessionId   string `json:"sessionId"`
	UserId      string `json:"userId"`
	CreatedTime int64  `json:"createdTime"`
	ExpiredTime int64  `json:"expiredTime"`
}
