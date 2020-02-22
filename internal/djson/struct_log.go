package djson

type ActionLog struct {
	Ip         string  `json:"ip"`
	OsGroup    OsGroup `json:"os_group"`
	SessionId  string  `json:"session_id"`
	CategoryId string  `json:"category_id"`
	EventApp   int     `json:"event_app"`
	EventId    string  `json:"event_id"`
	ArticleId  int     `json:"article_id"`

	TimeCreate int64 `json:"time_create"` //utc0 millisecond
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

type UserReport struct {
	TotalUser             int `json:"total_user"`
	NumNewUser            int `json:"num_new_user"`
	NumActivedUserToday   int `json:"num_actived_user_today"`
	NumActivedUserIn7Day  int `json:"num_actived_user_in_7_day"`
	NumActivedUserIn15Day int `json:"num_actived_user_in_15_day"`
}

type UsersLog struct {
	UserId          string `json:"user_id"`
	TimeCreate      int64  `json:"time_create"`
	LastUpdatedTime int64  `json:"last_updated_time,omitempty"`
}

type WebAction struct {
	Guid	string	`json:"guid"`
	Time_group TimeGroup `json:"time_group"`
	Ip         string  `json:"ip"`
	CategoryId string  `json:"category_id"`
	ArticleId  int     `json:"article_id"`
}

type TimeGroup struct {
	Cookie_create	int64	`json:"cookie_create"`
	Time_create		int64	`json:"time_create"`
}