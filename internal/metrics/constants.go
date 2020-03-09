package metrics

const (
	ReporWebLog = `stats.gauges.dora.web.log.%v`
	ReporCategoryWebLog = `stats.gauges.dora.web.log.category.%v.%v`
	RequestsSumMetric  = `stats.gauges.%v.dora.log.request.%v`
	ResponsesSumMetric = `stats.gauges.%v.dora.log.response.code.%v`
	StatusCodeMetric   = `stats.gauges.%v.dora.log.response.code.%v.%v`
	ConnectionMetric   = `stats.gauges.%v.dora.log.connection.%v`
	ConnectQuery       = `ss -nat | grep -w %v | awk {'print $1'} | cut -d':' -f1 | sort | uniq -c | sort -nr`

	GenericChannelMetric = `stats.gauges.%v.dora.log.latency.%v`
	MaxConnections       = 50000
)

var MatchingUrl = map[string]string{
	`/logging/trace`:     `trace-android`,
	`/logging/trace/dev`: `trace-ios`,
}
