package metrics

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kinghub-gateway/internal/utils"
	"github.com/marpaia/graphite-golang"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// DataChannel is the data you pass into the channel. Using Name we
// will put the Value into the right bucket.
type DataChannel struct {
	Name  string
	Value float64
}

type dataStore struct {
	sync.RWMutex
	data map[string][]float64
}

func NewDataStore() dataStore {
	return dataStore{data: make(map[string][]float64)}
}

func (ds dataStore) ResetKey(key string) {
	ds.Lock()
	defer ds.Unlock()
	ds.data[key] = make([]float64, 0)
}

func (ds dataStore) Get(key string) []float64 {
	ds.RLock()
	defer ds.RUnlock()
	return ds.data[key]
}

func (ds dataStore) Add(key string, value float64) {
	ds.data[key] = append(ds.data[key], value)
}

// GenericChannelAspect, exported fields are used to store json
// fields. All fields are measured in nanoseconds.
type GenericChannelAspect struct {
	gcdLock     sync.RWMutex
	name        string
	tempStore   dataStore
	internalGcd map[string]GenericChannelData
	graphite    *graphite.Graphite
	host        string
	port        int
}

// GenericChannelData
type GenericChannelData struct {
	Count     int       `json:"count"`
	Min       float64   `json:"min"`
	Max       float64   `json:"max"`
	Mean      float64   `json:"mean"`
	Stdev     float64   `json:"stdev"`
	P90       float64   `json:"p90"`
	P95       float64   `json:"p95"`
	P99       float64   `json:"p99"`
	Timestamp time.Time `json:"timestamp"`
}

// NewGenericChannelAspect returns a new initialized GenericChannelAspect
// object.
func NewGenericChannelAspect(name string, graphite *graphite.Graphite, host string, port int) *GenericChannelAspect {
	gc := &GenericChannelAspect{name: name}
	gc.tempStore = NewDataStore()
	gc.internalGcd = make(map[string]GenericChannelData)
	gc.graphite = graphite
	gc.host = host
	gc.port = port
	return gc
}

// StartTimer will call a forever loop in a goroutine to calculate
// metrics for measurements every d ticks.
func (gc *GenericChannelAspect) StartTimer(d time.Duration) {
	timer := time.Tick(d)
	go func() {
		for {
			<-timer
			gc.calculate()
		}
	}()
}

// SetupGenericChannelAspect returns an unbuffered channel for type
// DataChannel, such that you can send arbitrary key (string) value
// (float64) pairs to it.
func (gc *GenericChannelAspect) SetupGenericChannelAspect() chan DataChannel {
	lgc := gc // save gc in closure
	ch := make(chan DataChannel)
	go func() {
		for {
			lgc.add(<-ch)
		}
	}()
	return ch
}

// GetStats to fulfill aspects.Aspect interface, it returns a copy of
// the calculated data set that will be served as JSON.
func (gc *GenericChannelAspect) GetStats() interface{} {
	gc.gcdLock.RLock()
	defer gc.gcdLock.RUnlock()

	var mod bytes.Buffer
	enc := gob.NewEncoder(&mod)
	dec := gob.NewDecoder(&mod)

	err := enc.Encode(gc.internalGcd)
	if err != nil {
		return err
	}

	var cpy map[string]GenericChannelData
	err = dec.Decode(&cpy)
	if err != nil {
		return err
	}

	return cpy
}

// Name to fulfill aspects.Aspect interface, it will return the name
// of the JSON object that will be served.
func (gc *GenericChannelAspect) Name() string {
	return gc.name
}

// InRoot to fulfill aspects.Aspect interface, it will return where to
// put the JSON object into the monitoring endpoint.
func (gc *GenericChannelAspect) InRoot() bool {
	return false
}

func (gc *GenericChannelAspect) add(dc DataChannel) {
	gc.tempStore.Lock()
	defer gc.tempStore.Unlock()

	gc.tempStore.Add(dc.Name, dc.Value)
}

func (gc *GenericChannelAspect) calculate() {
	gc.tempStore.Lock()
	defer gc.tempStore.Unlock()
	for name, list := range gc.tempStore.data {
		sortedSlice := list[:]
		gc.tempStore.data[name] = make([]float64, 0)
		l := len(sortedSlice)

		// if tempStore is empty have to set everything to 0 and update timestamp
		if l < 1 {
			gc.gcdLock.Lock()
			gc.internalGcd[name] = GenericChannelData{Timestamp: time.Now()}
			gc.gcdLock.Unlock()
			continue
		}

		sort.Float64s(sortedSlice)
		m := mean(sortedSlice, l)

		gc.gcdLock.Lock()
		gc.internalGcd[name] = GenericChannelData{
			Timestamp: time.Now(),
			Count:     l,
			Min:       sortedSlice[0],
			Max:       sortedSlice[l-1],
			Mean:      m,
			Stdev:     correctedStdev(sortedSlice, m, l),
			P90:       p90(sortedSlice, l),
			P95:       p95(sortedSlice, l),
			P99:       p99(sortedSlice, l),
		}
		gc.gcdLock.Unlock()
	}
	//gc.gcdLock.Lock()
	//defer gc.gcdLock.Unlock()
	var timeTmp int64 = time.Now().Unix()
	//var mod bytes.Buffer
	//enc := gob.NewEncoder(&mod)
	//dec := gob.NewDecoder(&mod)
	//
	//err := enc.Encode(gc.internalGcd)
	//if err != nil {
	//	return
	//}
	//
	//var cpy map[string]GenericChannelData
	//err = dec.Decode(&cpy)
	//if err != nil {
	//	return
	//}
	cpy := gc.GetStats()
	switch cpy.(type) {
	case error:
		utils.HandleError(cpy)
	case map[string]GenericChannelData:
		go gc.Push(timeTmp, cpy.(map[string]GenericChannelData))
	default:
		fmt.Println(cpy)
	}
}

func GenericChannelHandler(gc *GenericChannelAspect) gin.HandlerFunc {
	genericCH := gc.SetupGenericChannelAspect()
	return func(c *gin.Context) {
		now := time.Now()
		c.Next()
		took := time.Since(now)
		genericCH <- DataChannel{Name: c.Request.URL.Path, Value: float64(took)}
	}
}

func (gc *GenericChannelAspect) Push(timeTmp int64, gcd map[string]GenericChannelData) {
	defer func() {
		if err := recover(); err != nil {
			utils.HandleError(err)
		}
	}()
	metrics := make([]graphite.Metric, 0)
	//connections
	numConnections := gc.GetConnections()
	for k, v := range numConnections {
		name := fmt.Sprintf(ConnectionMetric, gc.host, k)
		metrics = append(metrics, graphite.NewMetric(name, strconv.Itoa(v), timeTmp))
	}
	connEsTab := numConnections["estab"]
	ccu := 0
	for api, val := range gcd {
		api = strings.Replace(api, ".", "_", -1)
		api = strings.Replace(api, "/", "_", -1)
		//p90
		name := fmt.Sprintf(GenericChannelMetric, gc.host, api+`.p90`)
		v := (int)(val.P90 / 1000)
		metrics = append(metrics, graphite.NewMetric(name, strconv.Itoa(v), timeTmp))
		//p95
		name = fmt.Sprintf(GenericChannelMetric, gc.host, api+`.p95`)
		v = (int)(val.P95 / 1000)
		metrics = append(metrics, graphite.NewMetric(name, strconv.Itoa(v), timeTmp))
		//p99
		name = fmt.Sprintf(GenericChannelMetric, gc.host, api+`.p99`)
		v = (int)(val.P99 / 1000)
		metrics = append(metrics, graphite.NewMetric(name, strconv.Itoa(v), timeTmp))
		//Min
		name = fmt.Sprintf(GenericChannelMetric, gc.host, api+`.min`)
		v = (int)(val.Min / 1000)
		metrics = append(metrics, graphite.NewMetric(name, strconv.Itoa(v), timeTmp))
		//Max
		name = fmt.Sprintf(GenericChannelMetric, gc.host, api+`.max`)
		v = (int)(val.Max / 1000)
		metrics = append(metrics, graphite.NewMetric(name, strconv.Itoa(v), timeTmp))
		//Mean
		name = fmt.Sprintf(GenericChannelMetric, gc.host, api+`.mean`)
		v = (int)(val.Mean / 1000)
		metrics = append(metrics, graphite.NewMetric(name, strconv.Itoa(v), timeTmp))
		if v == 0 {
			ccu = MaxConnections
		} else if v > 1000000 {
			ccu = MaxConnections
		} else {
			ccu = MaxConnections * (1000000 / v)
		}
		name = fmt.Sprintf(ConnectionMetric, gc.host, `ccu`)
		metrics = append(metrics, graphite.NewMetric(name, strconv.Itoa(connEsTab), timeTmp))
		name = fmt.Sprintf(ConnectionMetric, gc.host, `max-ccu`)
		metrics = append(metrics, graphite.NewMetric(name, strconv.Itoa(ccu), timeTmp))
		//Stdev
		name = fmt.Sprintf(GenericChannelMetric, gc.host, api+`.stdev`)
		v = (int)(val.Stdev / 1000)
		metrics = append(metrics, graphite.NewMetric(name, strconv.Itoa(v), timeTmp))

	}
	if len(metrics) > 0 {
		_ = gc.graphite.SendMetrics(metrics)
	}
}

func (gc *GenericChannelAspect) GetConnections() map[string]int {
	cmd := fmt.Sprintf(ConnectQuery, gc.port)
	out, err := exec.Command("bash", "-c", cmd).Output()
	outData := make(map[string]int)
	outData["time-wait"] = 0
	outData["estab"] = 0
	if err != nil {
		return outData
	}
	data := string(out)
	tmpdata := strings.Split(data, "\n")
	if len(tmpdata) > 0 {
		for _, tmp := range tmpdata {
			var re = regexp.MustCompile(`[ ]{2,}`)
			tmp = re.ReplaceAllString(tmp, ` `)
			values := strings.Split(strings.Trim(tmp, " "), " ")
			if len(values) == 2 {
				numConnection, _ := strconv.Atoi(strings.Trim(values[0], " "))
				outData[strings.ToLower(values[1])] = numConnection
			}
		}
	}
	return outData
}
