package main

import (
	"context"
	"log"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	oltpmetric "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/metric"
	metricsdk "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

var res = resource.NewWithAttributes(
	semconv.SchemaURL,
	semconv.ServiceName("runtime-instrumentation-example"),
)

func main() {
	stdoutExp, err := stdoutmetric.New()
	if err != nil {
		log.Fatal(err)
	}
	// Register the exporter with an SDK via a periodic reader.
	stdoutRead := metricsdk.NewPeriodicReader(stdoutExp, metricsdk.WithInterval(1*time.Second))

	secureOpt := oltpmetric.WithInsecure()

	//read := metricsdk.NewManualReader()
	provider := metricsdk.NewMeterProvider(metricsdk.WithResource(res), metricsdk.WithReader(stdoutRead))
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
}
