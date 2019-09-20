package metrics

import (
//"fmt"
//"github.com/gin-gonic/gin"
////"github.com/Dora-Logs/internal/client"
//"github.com/Dora-Logs/internal/utils"
//"github.com/marpaia/graphite-golang"
//"strconv"
//"strings"
//"sync"
//"time"
)

// CounterHandler is a Gin middleware function that increments a
// global counter on each request.
//func CounterHandler(ca *CounterAspect) gin.HandlerFunc {
//	return func(ctx *gin.Context) {
//		ctx.Next()
//		ca.inc <- tuple{
//			path: ctx.Request.URL.Path,
//			code: ctx.Writer.Status(),
//		}
//	}
//}
//
//type tuple struct {
//	path string
//	code int
//}
//
//// CounterAspect stores a counter
//type CounterAspect struct {
//	services map[string][]client.Service
//
//	inc                     chan tuple
//	internalRequestsSum     int
//	internalRequests        map[string]int
//	internalRequestCodes    map[string]map[int]int
//	internalResponseCodeSum map[int]int
//
//	RequestsSum     int                    `json:"request_sum_per_minute"`
//	Requests        map[string]int         `json:"requests_per_minute"`
//	RequestCodes    map[string]map[int]int `json:"request_codes_per_minute"`
//	ResponseCodeSum map[int]int            `json:"request_sum_per_minute"`
//
//	graphite *graphite.Graphite
//	host     string
//	port     int
//
//	counterLock sync.RWMutex
//}
//
//// NewCounterAspect returns a new initialized CounterAspect object.
//func NewCounterAspect(services map[string][]client.Service, graphite *graphite.Graphite, host string) *CounterAspect {
//	ca := &CounterAspect{}
//	ca.services = services
//	ca.inc = make(chan tuple)
//	ca.internalRequestsSum = 0
//	ca.internalRequests = make(map[string]int)
//	ca.internalRequestCodes = make(map[string]map[int]int)
//	ca.internalResponseCodeSum = make(map[int]int)
//
//	ca.graphite = graphite
//	ca.host = host
//	ca.counterLock = sync.RWMutex{}
//	return ca
//}
//
//// StartTimer will call a forever loop in a goroutine to calculate
//// metrics for measurements every d ticks. The parameter of this
//// function should normally be 1 * time.Minute, if not it will expose
//// unintuive JSON keys (requests_per_minute and
//// request_sum_per_minute).
//func (ca *CounterAspect) StartTimer(d time.Duration) {
//	timer := time.Tick(d)
//	go func() {
//		for {
//			select {
//			case tup := <-ca.inc:
//				ca.increment(tup)
//			case <-timer:
//				ca.reset()
//			}
//		}
//	}()
//}
//
//// GetStats to fulfill aspects.Aspect interface, it returns the data
//// that will be served as JSON.
//func (ca *CounterAspect) GetStats() interface{} {
//	return *ca
//}
//
//// Name to fulfill aspects.Aspect interface, it will return the name
//// of the JSON object that will be served.
//func (ca *CounterAspect) Name() string {
//	return "Counter"
//}
//
//// InRoot to fulfill aspects.Aspect interface, it will return where to
//// put the JSON object into the monitoring endpoint.
//func (ca *CounterAspect) InRoot() bool {
//	return false
//}
//
//func (ca *CounterAspect) increment(tup tuple) {
//	ca.counterLock.Lock()
//	ca.internalRequestsSum++
//	ca.internalRequests[tup.path]++
//	ca.internalResponseCodeSum[tup.code]++
//
//	if _, ok := ca.internalRequestCodes[tup.path]; !ok {
//		ca.internalRequestCodes[tup.path] = make(map[int]int)
//	}
//	ca.internalRequestCodes[tup.path][tup.code]++
//	ca.counterLock.Unlock()
//}
//
//func (ca *CounterAspect) reset() {
//	ca.counterLock.Lock()
//	ca.RequestsSum = ca.internalRequestsSum
//	ca.Requests = ca.internalRequests
//	ca.RequestCodes = ca.internalRequestCodes
//	ca.ResponseCodeSum = ca.internalResponseCodeSum
//
//	ca.internalRequestsSum = 0
//	ca.internalRequests = make(map[string]int, ca.RequestsSum)
//	ca.internalRequestCodes = make(map[string]map[int]int, len(ca.RequestCodes))
//	ca.internalResponseCodeSum = make(map[int]int)
//
//	ca.counterLock.Unlock()
//	var timeTmp int64 = time.Now().Unix()
//	go ca.Push(timeTmp, ca.RequestsSum, ca.Requests, ca.RequestCodes, ca.ResponseCodeSum)
//}
//
//func (ca *CounterAspect) Push(timeTmp int64, RequestsSum int, Requests map[string]int, RequestCodes map[string]map[int]int, ResponseCodeSum map[int]int) {
//
//	defer func() {
//		if err := recover(); err != nil {
//			utils.HandleError(err)
//		}
//	}()
//
//	metrics := make([]graphite.Metric, 0)
//
//	//total request to host
//	total := fmt.Sprintf(RequestsSumMetric, ca.host, `total-endpoint`)
//	metrics = append(metrics, graphite.NewMetric(total, strconv.Itoa(RequestsSum), timeTmp))
//
//	//total respose code each host
//	for code, val := range ResponseCodeSum {
//		name := fmt.Sprintf(ResponsesSumMetric, ca.host, code)
//		metrics = append(metrics, graphite.NewMetric(name, strconv.Itoa(val), timeTmp))
//	}
//
//	// request each api
//	for api, val := range Requests {
//		fromPath := strings.SplitN(api, "/", 4)
//		if len(fromPath) > 3 && strings.EqualFold(fromPath[1], "api") {
//			fromPathMain := fromPath[3]
//			fromPathMain = fmt.Sprintf("/%s", fromPathMain)
//			if _, ok := ca.services[fromPathMain]; ok {
//				api = strings.Replace(api, ".", "_", -1)
//				api = strings.Replace(api, "/", "_", -1)
//				name := fmt.Sprintf(RequestsSumMetric, ca.host, api)
//				metrics = append(metrics, graphite.NewMetric(name, strconv.Itoa(val), timeTmp))
//			}
//		}
//	}
//
//	//code for each api
//	for api, codes := range RequestCodes {
//		fromPath := strings.SplitN(api, "/", 4)
//		if len(fromPath) > 3 && strings.EqualFold(fromPath[1], "api") {
//			fromPathMain := fromPath[3]
//			fromPathMain = fmt.Sprintf("/%s", fromPathMain)
//			if _, ok := ca.services[fromPathMain]; ok {
//				api = strings.Replace(api, ".", "_", -1)
//				api = strings.Replace(api, "/", "_", -1)
//				for code, val := range codes {
//					name := fmt.Sprintf(StatusCodeMetric, ca.host, api, code)
//					metrics = append(metrics, graphite.NewMetric(name, strconv.Itoa(val), timeTmp))
//				}
//			}
//		}
//	}
//
//	if len(metrics) > 0 {
//		err := ca.graphite.SendMetrics(metrics)
//		if err != nil {
//			utils.HandleError(err)
//		}
//	}
//}
