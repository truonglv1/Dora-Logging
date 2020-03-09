package metrics

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/Dora-Logging/internal/djson"
	"github.com/Dora-Logging/utils"
	"github.com/gin-gonic/gin"
	"github.com/marpaia/graphite-golang"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

// CounterHandler is a Gin middleware function that increments a
// global counter on each request.
func CounterHandler(ca *CounterAspect) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
		ca.inc <- tuple{
			path: ctx.Request.URL.Path,
			code: ctx.Writer.Status(),
		}
	}
}

type tuple struct {
	path string
	code int
}

// CounterAspect stores a counter
type CounterAspect struct {
	inc                  chan tuple
	internalRequestsSum  int
	internalRequests     map[string]int
	internalRequestCodes map[string]map[int]int

	internalTotalUser	 int
	internalDAU			 int

	RequestsSum          int                    `json:"request_sum_per_minute"`
	Requests             map[string]int         `json:"requests_per_minute"`
	RequestCodes         map[string]map[int]int `json:"request_codes_per_minute"`

	TotalUser			 int 					`json:"total_user"`
	DailyActiveUser      int                     `json:"daily_active_user"`

	categories 			 map[string]string
	graphite             *graphite.Graphite
	host                 string
	counterLock          sync.RWMutex
}

// NewCounterAspect returns a new initialized CounterAspect object.
func NewCounterAspect(graphite *graphite.Graphite, host string, categories map[string]string) *CounterAspect {
	ca := &CounterAspect{}
	ca.inc = make(chan tuple)

	ca.internalRequestsSum = 0
	ca.internalRequests = make(map[string]int)
	ca.internalRequestCodes = make(map[string]map[int]int)

	ca.internalTotalUser=0
	ca.internalDAU=0

	ca.categories = categories
	ca.graphite = graphite
	ca.host = host
	ca.counterLock = sync.RWMutex{}
	return ca
}

// StartTimer will call a forever loop in a goroutine to calculate
// metrics for measurements every d ticks. The parameter of this
// function should normally be 1 * time.Minute, if not it will expose
// unintuive JSON keys (requests_per_minute and
// request_sum_per_minute).
func (ca *CounterAspect) StartTimer(d time.Duration) {
	timer := time.NewTicker(d).C
	go func() {
		for {
			select {
			case tup := <-ca.inc:
				ca.increment(tup)
			case <-timer:
				ca.reset()
			}
		}
	}()
}

//// GetStats to fulfill aspects.Aspect interface, it returns the data
//// that will be served as JSON.
//func (ca *CounterAspect) GetStats() interface{} {
//	return *ca
//}

// Name to fulfill aspects.Aspect interface, it will return the name
// of the JSON object that will be served.
func (ca *CounterAspect) Name() string {
	return "Counter"
}

// InRoot to fulfill aspects.Aspect interface, it will return where to
// put the JSON object into the monitoring endpoint.
func (ca *CounterAspect) InRoot() bool {
	return false
}

func (ca *CounterAspect) increment(tup tuple) {
	ca.counterLock.Lock()
	ca.internalRequestsSum++
	ca.internalRequests[tup.path]++
	if _, ok := ca.internalRequestCodes[tup.path]; !ok {
		ca.internalRequestCodes[tup.path] = make(map[int]int)
	}
	ca.internalRequestCodes[tup.path][tup.code]++
	ca.counterLock.Unlock()
}

func (ca *CounterAspect) reset() {
	ca.counterLock.Lock()
	ca.RequestsSum = ca.internalRequestsSum
	ca.Requests = ca.internalRequests
	ca.RequestCodes = ca.internalRequestCodes

	ca.internalRequestsSum = 0
	ca.internalRequests = make(map[string]int, ca.RequestsSum)
	ca.internalRequestCodes = make(map[string]map[int]int, len(ca.RequestCodes))
	ca.counterLock.Unlock()
	var timeTmp = time.Now().Unix()
	go ca.Push(timeTmp, ca.RequestsSum, ca.Requests, ca.RequestCodes)
}

func (ca *CounterAspect) Push(timeTmp int64, RequestsSum int, Requests map[string]int, RequestCodes map[string]map[int]int) {
	defer func() {
		if err := recover(); err != nil {
			utils.HandleError(err)
		}
	}()

	metrics := make([]graphite.Metric, 0)
	//total
	total := fmt.Sprintf(RequestsSumMetric, ca.host, `total`)
	metrics = append(metrics, graphite.NewMetric(total, strconv.Itoa(RequestsSum), timeTmp))

	//total user
	totalUser:=getTotalUser(1)
	totalUserWeb := fmt.Sprintf(ReporWebLog, `total_user`)
	metrics = append(metrics, graphite.NewMetric(totalUserWeb, strconv.Itoa(totalUser), timeTmp))
	//total DAU
	dau, reportCate := ca.getDailyActiveUser()
	totalDAU := fmt.Sprintf(ReporWebLog, `total_dau`)
	metrics = append(metrics, graphite.NewMetric(totalDAU, strconv.Itoa(dau), timeTmp))

	//report cate
	for key, val := range reportCate{
		//fmt.Println(key, "_", val)
		totalUserViewCate := fmt.Sprintf(ReporCategoryWebLog, key, `total_user_view`)
		metrics = append(metrics, graphite.NewMetric(totalUserViewCate, strconv.Itoa(val), timeTmp))

	}

	//api
	for api, val := range Requests {
		if tmp, ok := MatchingUrl[api]; ok {
			name := fmt.Sprintf(RequestsSumMetric, ca.host, tmp)
			metrics = append(metrics, graphite.NewMetric(name, strconv.Itoa(val), timeTmp))
		}
	}

	for api, name := range MatchingUrl {
		if _, ok := Requests[api]; !ok {
			name := fmt.Sprintf(RequestsSumMetric, ca.host, name)
			metrics = append(metrics, graphite.NewMetric(name, "0", timeTmp))
		}
	}
	// status code
	for key, api := range MatchingUrl {
		codes := make(map[int]int)
		if _, ok := RequestCodes[key]; ok {
			codes = RequestCodes[key]
		}
		//fake data
		if _, ok := codes[http.StatusOK]; !ok {
			codes[http.StatusOK] = 0
		}
		if _, ok := codes[http.StatusRequestTimeout]; !ok {
			codes[http.StatusRequestTimeout] = 0
		}
		if _, ok := codes[http.StatusInternalServerError]; !ok {
			codes[http.StatusInternalServerError] = 0
		}
		if _, ok := codes[http.StatusBadGateway]; !ok {
			codes[http.StatusBadGateway] = 0
		}
		if _, ok := codes[http.StatusServiceUnavailable]; !ok {
			codes[http.StatusServiceUnavailable] = 0
		}

		for code, val := range codes {
			name := fmt.Sprintf(StatusCodeMetric, ca.host, api, code)
			metrics = append(metrics, graphite.NewMetric(name, strconv.Itoa(val), timeTmp))
		}
	}

	//send metrics
	if len(metrics) > 0 {
		err := ca.graphite.SendMetrics(metrics)
		if err != nil {
			utils.HandleError(err)
		}
	}
}

func getTotalUser(numday int) int  {
	totalUser :=0

	userMapOld := make(map[string]djson.WebAction)
	userMap := make(map[string]djson.WebAction)
	//read file
	oldFile, err := os.Open("report/users.log")
	if err != nil {
		utils.HandleError(err)
	}
	report := bufio.NewScanner(oldFile)
	for report.Scan(){
		userMapOld[report.Text()] = djson.WebAction{}
	}

	for i:=0;i<numday;i++{
		var path string
		if i==0 {
			path = "logging/web-log.log"
		}else {
			path = "logging/web-log.log."+time.Now().AddDate(0, 0, -i).Format("2006-01-02");
		}
		file, err := os.Open(path)
		if err != nil {
			utils.HandleError(err)
		}
		logging := bufio.NewScanner(file)
		for logging.Scan(){
			var w djson.WebAction
			if err := json.Unmarshal(logging.Bytes(), &w); err != nil {
				utils.HandleError(err)
			}
			_,ok := userMapOld[w.Guid]
			if !ok{
				userMap[w.Guid] = w
			}
		}
	}
	f, err := os.OpenFile("report/users.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	for key, _ := range userMap {
		if _, err := f.WriteString(string(fmt.Sprintf("\n%s",key))); err != nil {
			log.Println(err)
		}
	}
	totalUser = len(userMap) + len(userMapOld)
	return totalUser
}

func (ca *CounterAspect) getDailyActiveUser() (int,map[string]int) {
	counterReport := make(map[string]int)

	cate := make(map[string]map[string]string)
	userMap := make(map[string]djson.WebAction)

	file, err := os.Open("logging/web-log.log")
	if err != nil {
		utils.HandleError(err)
	}
	logging := bufio.NewScanner(file)
	for logging.Scan(){
		var w djson.WebAction
		if err := json.Unmarshal(logging.Bytes(), &w); err != nil {
			utils.HandleError(err)
		}
		_,existUser := userMap[w.Guid]
		if !existUser{
			userMap[w.Guid] = w
		}
		_, existCate := cate[w.CategoryId]
		if existCate{
			_,existUserInCate := cate[w.CategoryId][w.Guid]
			if !existUserInCate{
				cate[w.CategoryId][w.Guid] = w.Guid
			}
		}else {
			cate[w.CategoryId] = make(map[string]string)
			cate[w.CategoryId][w.Guid] = w.Guid
		}
	}

	//report category (total user view category)
	for key, val := range cate{
		_,ok := ca.categories[key]
		if ok{
			counterReport[ca.categories[key]] = len(val)
		}
	}

	return len(userMap), counterReport
}

