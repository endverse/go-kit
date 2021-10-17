package opentracing

import (
	"io"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	jaegerconfig "github.com/uber/jaeger-client-go/config"
)

func InitializeTracing(serviceName, tracingHostPort, tracingUser, tracingPassword string) (closer io.Closer, err error) {
	return createGlobalJager(&jaegerconfig.Configuration{
		ServiceName: serviceName,
		Disabled:    false,
		Sampler: &jaegerconfig.SamplerConfig{
			Type: jaeger.SamplerTypeConst,
			// The param's value is between 0 and 1,
			// if set to 1, it will output all operations to the Reporter.
			Param: 1,
		},
		Reporter: &jaegerconfig.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: tracingHostPort,
			User:               tracingUser,
			Password:           tracingPassword,
		},
	})
}

// createGlobalJager init a Tracer and return a Closer
func createGlobalJager(cfg *config.Configuration) (closer io.Closer, err error) {
	tracer, closer, err := (cfg).NewTracer()
	if err != nil {
		return
	}

	// Set global tracer, if not set,
	// it can not generate a span with contxt.
	opentracing.SetGlobalTracer(tracer)
	return
}
