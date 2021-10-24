package app

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/bhoriuchi/go-bunyan/bunyan"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	AppConfig "github.com/webalytic.go/cmd/collector/config"
	CommonCfg "github.com/webalytic.go/common/config"
	Datasources "github.com/webalytic.go/common/datasources"
	RedisBroker "github.com/webalytic.go/common/redis"
)

var (
	collectRequestCnt = prometheus.NewCounter(prometheus.CounterOpts{
		Name:        "collector_collect_request_total",
		Subsystem:   "Collector",
		Help:        "Number of collected POST request",
		ConstLabels: prometheus.Labels{"version": "1"},
	})
	collectRequestInProgress = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "collector_collect_request_in_progress",
		Help:      "Number of in progress collect requests",
		Subsystem: "Collector",
	})
	errorNegRequestCnt = prometheus.NewCounter(prometheus.CounterOpts{
		Name:        "collector_collect_request_failed",
		Subsystem:   "Collector",
		Help:        "Number of incorrect collect POST requests",
		ConstLabels: prometheus.Labels{"version": "1"},
	})
	collectRequesHandleLatency = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:      "collector_collect_request_handling_latency_ms",
		Subsystem: "Collector",
		Help:      "Request latency histogram",
		Buckets:   []float64{1, 5, 25, 50, 100, 250, 500, 1000, math.Inf(+1)},
	})
)

func getTrace(logger bunyan.Logger) string {
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		logger.Fatal(fmt.Sprintf("Unable to generate traceID with error: %s", err))
		return ""
	}
	return hex.EncodeToString(bytes)
}

func InitCollectHandler() {
	prometheus.MustRegister(collectRequestCnt)
	prometheus.MustRegister(errorNegRequestCnt)
	prometheus.MustRegister(collectRequestInProgress)
	prometheus.MustRegister(collectRequesHandleLatency)
}

type ICollectHandler struct {
	Handler http.HandlerFunc
}

func CollectHandler(
	logger bunyan.Logger,
	broker *RedisBroker.RedisBroker,
	cfg *AppConfig.CollectorConfig,
	redisCfg *CommonCfg.RedisConfig) *ICollectHandler {
	streamName := redisCfg.StreamName()

	return &ICollectHandler{
		Handler: func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()
			collectRequestCnt.Inc()
			collectRequestInProgress.Inc()
			logger.Debug("Collect handler")
			vars := mux.Vars(r)
			logger.Debug(fmt.Sprintf("Request vars: %s", vars["category"]))

			if r.Body == nil {
				errorNegRequestCnt.Inc()
				logger.Error("Empty body is not allowed")
				http.Error(w, "Empty body is not allowed", http.StatusBadRequest)
				collectRequestInProgress.Dec()
				return
			}
			traceID := getTrace(logger)
			var out []byte
			var err error
			switch vars["category"] {
			case "payment":
				out, err = Datasources.PaymentHelper(r.Body, traceID, logger)
			default:
				err = nil
				out = nil
				http.Error(w, "Empty body is not allowed", http.StatusBadRequest)
				return
			}
			if err != nil {
				errorNegRequestCnt.Inc()
				collectRequestInProgress.Dec()
				http.Error(w, err.Error(), http.StatusBadRequest)
			}

			logger.Debug("traceID: %s: Endpoint collect data: %s, type: %s, channel: %s", traceID, vars["category"], string(out), streamName)
			broker.Publish(streamName, out)

			collectRequestInProgress.Dec()
			collectRequesHandleLatency.Observe(
				float64(time.Since(startTime).Milliseconds()),
			)
			fmt.Fprintf(w, "OK")
		},
	}
}
