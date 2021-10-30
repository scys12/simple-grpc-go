package tracer

import (
	"context"
	"io"

	"github.com/opentracing/opentracing-go"
	"github.com/scys12/simple-grpc-go/pkg/logger"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"
)

func Init(appName string) (io.Closer, error) {
	cfg := config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: "jaeger:6831",
		},
	}

	jLogger := jaeger.NullLogger
	jMetricsFactory := metrics.NullFactory

	closer, err := cfg.InitGlobalTracer(
		appName,
		config.Logger(jLogger),
		config.Metrics(jMetricsFactory),
	)
	if err != nil {
		logger.Log.Error("[Tracer] could not initialize jaeger, err : %v", zap.String("jaeger", err.Error()))
		return nil, err
	}
	return closer, nil
}

func StartSpanFromContext(ctx context.Context, operationName string) (opentracing.Span, context.Context) {
	return opentracing.StartSpanFromContext(ctx, operationName)
}
