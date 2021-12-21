package main

import (
	"fmt"
	"math"
	"net/http"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func primeNumbers() []int {
	var primes []int

	for i := 2; i < 10000; i++ {
		isPrime := true

		for j := 2; j <= int(math.Sqrt(float64(i))); j++ {
			if i%j == 0 {
				isPrime = false
				break
			}
		}

		if isPrime {
			primes = append(primes, i)
		}
	}
	return primes
}

func PrimeNumbersBenchmark(N int) {
	for i := 0; i < N; i++ {
		_ = primeNumbers()
	}
}

func recordMetrics() {

	go func() {
		for {
			start := time.Now()
			N := 500
			PrimeNumbersBenchmark(N)
			duration := time.Since(start)
			opsPerSecond := float64(N) / duration.Seconds()

			fmt.Printf("Test finished in %v ms (%0.0f op/s)\n", duration.Milliseconds(), opsPerSecond)
			perfTestOps.Set(opsPerSecond)
			time.Sleep(300 * time.Second)
		}
	}()
}

var (
	perfTestOps = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "cpu_performance_test_ops",
		Help: "The number of operations per second the test executed.",
	})
)

func main() {
	testing.Init()
	recordMetrics()

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}
