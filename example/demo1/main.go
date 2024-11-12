package main

import (
	"fmt"
	"github.com/liangweijiang/go-metric/meter"
	"github.com/liangweijiang/go-metric/pkg/config"
	"net/http"
)

func main() {
	m, err := meter.NewMeter(
		meter.WithProviderType(config.MeterProviderTypePrometheus),
		meter.WithEnv(config.MeterEnvTest),
		meter.WithRuntimeMetricsCollector(),
		meter.WithPrometheusPort(16666))
	if err != nil {
		fmt.Println(err)
	}
	meter.SetGlobalMeter(m)

	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		m.NewCounter("http_counter", "", "").AddTag("test1", "test1").IncrOne(r.Context())
		fmt.Fprintf(w, "Hello, World!")
	})

	http.HandleFunc("/closeMeter", func(w http.ResponseWriter, r *http.Request) {
		m.WithRunning(false)
	})

	http.HandleFunc("/startMeter", func(w http.ResponseWriter, r *http.Request) {
		m.WithRunning(true)
	})

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		h := m.GetHandler()
		if h == nil {
			fmt.Println("handler is nil")
			return
		}
		h.ServeHTTP(w, r)

	})
	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
