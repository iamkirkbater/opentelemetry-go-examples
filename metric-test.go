    //
    // 	read := metric.NewPeriodicReader(exp, metric.WithInterval(1*time.Second))
    // 	provider := metric.NewMeterProvider(metric.WithResource(res), metric.WithReader(read))
    // 	defer func() {
    // 		err := provider.Shutdown(context.Background())
    // 		if err != nil {
    // 			log.Fatal(err)
    // 		}
    // 	}()
    // 	otel.SetMeterProvider(provider)
    // 	m := provider.Meter("meter-name")
    //
    // 	log.Print("registering metric")
    //
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
        "go.opentelemetry.io/otel/attribute"
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
        exp, err := stdoutmetric.New()
        if err != nil {
            log.Fatal(err)
        }

        // Register the exporter with an SDK via a periodic reader.
        read := metricsdk.NewPeriodicReader(exp, metricsdk.WithInterval(1*time.Second))
        provider := metricsdk.NewMeterProvider(metricsdk.WithResource(res), metricsdk.WithReader(read))
        defer func() {
            err := provider.Shutdown(context.Background())
            if err != nil {
                log.Fatal(err)
            }
        }()
        otel.SetMeterProvider(provider)

        ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
        defer cancel()

        log.Print("Starting runtime instrumentation:")
        m := otel.Meter("my.meter.name")

        counter, _ := m.Int64Counter(
            "some.prefix.counter",
            metric.WithDescription("my-counter"),
            metric.WithUnit("calls"),
        )

        counter.Add(ctx, 1, metric.WithAttributes(
			attribute.String("cmd", "root")),
		)

        <-ctx.Done()
    }
