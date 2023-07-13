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

func main() {
	exp, err := stdoutmetric.New()
	if err != nil {
		log.Fatal(err)
	}

	// Register the exporter with an SDK via a periodic reader.
	//read := metricsdk.NewPeriodicReader(exp, metricsdk.WithInterval(1*time.Second))
	read := metricsdk.NewManualReader()
	provider := metricsdk.NewMeterProvider(metricsdk.WithResource(res), metricsdk.WithReader(read))
	defer func() {
		err := provider.Shutdown(context.Background())
		if err != nil {
			log.Fatal(err)
		}
	}()
	otel.SetMeterProvider(provider)

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

	collectedMetrics := &metricdata.ResourceMetrics{}
	read.Collect(context.TODO(), collectedMetrics)
	exp.Export(context.TODO(), collectedMetrics)
}
