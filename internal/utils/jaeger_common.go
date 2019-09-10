package utils

import (
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"io"
)

func InitJaeger(serviceName string) (opentracing.Tracer, io.Closer) {

	cfg := jaegercfg.Configuration{
		ServiceName: serviceName,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  "probabilistic",
			Param: 0.001,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LocalAgentHostPort: "10.3.68.12:6831",
			LogSpans:           true,
		},
	}

	tracer, closer, err := cfg.NewTracer(jaegercfg.Logger(jaeger.NullLogger))

	if err != nil {
		panic(fmt.Sprintf("Could not initialize jaeger tracer: %s", err.Error()))
	}

	return tracer, closer
}
