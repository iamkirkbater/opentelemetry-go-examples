package main

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/metric"
	metricsdk "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

var res = resource.NewWithAttributes(
	semconv.SchemaURL,
	semconv.ServiceName("runtime-instrumentation-example"),
)

var metricReader = metricsdk.NewManualReader()

func setup_metrics() func() {
	ctx := context.TODO()
	provider := metricsdk.NewMeterProvider(metricsdk.WithResource(res), metricsdk.WithReader(metricReader))
	otel.SetMeterProvider(provider)

	return func() {
		exp, err := stdoutmetric.New()
		if err != nil {
			log.Fatal(err)
		}

		collectedMetrics := &metricdata.ResourceMetrics{}
		metricReader.Collect(ctx, collectedMetrics)
		exp.Export(ctx, collectedMetrics)

		err = provider.Shutdown(context.Background())
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	shutdown := setup_metrics()
	defer shutdown()

	log.Print("Starting runtime instrumentation:")
	m := otel.Meter("my.meter.name")

	counter, _ := m.Int64Counter(
		"some.prefix.counter",
		metric.WithDescription("my-counter"),
		metric.WithUnit("calls"),
	)

	counter.Add(context.TODO(), 1, metric.WithAttributes(
		attribute.String("cmd", "root")),
	)
}
