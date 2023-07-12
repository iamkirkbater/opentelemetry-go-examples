// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build go1.18
// +build go1.18

package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"go.opentelemetry.io/otel"
	metricExporter "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	metricOpts "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

var res = resource.NewWithAttributes(
	semconv.SchemaURL,
	semconv.ServiceName("occ"),
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	exp, err := metricExporter.New(ctx, metricExporter.WithInsecure(), metricExporter.WithEndpoint("localhost:4318"))
	if err != nil {
		log.Fatal(err)
	}

	read := metric.NewPeriodicReader(exp, metric.WithInterval(1*time.Second))
	provider := metric.NewMeterProvider(metric.WithResource(res), metric.WithReader(read))
	defer func() {
		err := provider.Shutdown(context.Background())
		if err != nil {
			log.Fatal(err)
		}
	}()
	otel.SetMeterProvider(provider)
	m := provider.Meter("meter-name")

	log.Print("registering metric")
	counter, _ := m.Int64Counter(
		"some.prefix.counter",
		metricOpts.WithDescription("test description"),
	)

	counter.Add(ctx, 1)
}
