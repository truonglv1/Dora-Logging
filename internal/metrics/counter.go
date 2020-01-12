package metrics

import (
	"fmt"
	"github.com/Dora-Logging/utils"
	"github.com/gin-gonic/gin"
	"github.com/marpaia/graphite-golang"
	"net/http"
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
	RequestsSum          int                    `json:"request_sum_per_minute"`
	Requests             map[string]int         `json:"requests_per_minute"`
	RequestCodes         map[string]map[int]int `json:"request_codes_per_minute"`
	graphite             *graphite.Graphite
	host                 string
	counterLock          sync.RWMutex
}

// NewCounterAspect returns a new initialized CounterAspect object.
func NewCounterAspect(graphite *graphite.Graphite, host string) *CounterAspect {
	ca := &CounterAspect{}
	ca.inc = make(chan tuple)
	ca.internalRequestsSum = 0
	ca.internalRequests = make(map[string]int)
	ca.internalRequestCodes = make(map[string]map[int]int)
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
