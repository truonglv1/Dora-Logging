package metrics

const (
	RequestsSumMetric  = `stats.gauges.%v.kinghub.gateway.request.%v`
	ResponsesSumMetric = `stats.gauges.%v.kinghub.gateway.response.code.%v`
	StatusCodeMetric   = `stats.gauges.%v.kinghub.gateway.response.code.%v.%v`
	ConnectionMetric   = `stats.gauges.%v.kinghub.gateway.connection.%v`
	ConnectQuery       = `ss -nat | grep -w %v | awk {'print $1'} | cut -d':' -f1 | sort | uniq -c | sort -nr`

	GenericChannelMetric = `stats.gauges.%v.kinghub.gateway.latency.%v`
	MaxConnections       = 50000
)
