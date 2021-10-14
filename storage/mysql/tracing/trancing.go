package tracing

import (
	"io"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
)

// CreateGlobalJager init a Tracer and return a Closer
func CreateGlobalJager(cfg *config.Configuration) (closer io.Closer, err error) {
	tracer, closer, err := (cfg).NewTracer()
	if err != nil {
		return
	}

	// Set global tracer, if not set,
	// it can not generate a span with contxt.
	opentracing.SetGlobalTracer(tracer)
	return
}
