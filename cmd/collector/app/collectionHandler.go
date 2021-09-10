package app

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bhoriuchi/go-bunyan/bunyan"
	AppConfig "github.com/webalytic.go/cmd/collector/config"
	CommonCfg "github.com/webalytic.go/common/config"
	Datasources "github.com/webalytic.go/common/datasources"
	RedisBroker "github.com/webalytic.go/common/redis"
)

func getTrace(logger bunyan.Logger) string {
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		logger.Fatal(fmt.Sprintf("Unable to generate traceID with error: %s", err))
		return ""
	}
	return hex.EncodeToString(bytes)
}

func CollectHandler(
	logger bunyan.Logger,
	broker *RedisBroker.RedisBroker,
	cfg *AppConfig.CollectorConfig,
	redisCfg *CommonCfg.RedisConfig) http.HandlerFunc {
	streamName := redisCfg.StreamName()
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Debug("Collect handler")
		var payment Datasources.Payment
		if r.Body == nil {
			logger.Error("Empty body is not allowed")
			http.Error(w, "Empty body is not allowed", http.StatusBadRequest)
		}
		err := json.NewDecoder(r.Body).Decode(&payment)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		payment.TraceID = getTrace(logger)
		out, err := json.Marshal(payment)
		if err != nil {
			logger.Error(err)
		}
		logger.Debug("traceID: %s: Endpoint collect payment: %s, channel: %s", payment.TraceID, string(out), streamName)
		broker.Publish(streamName, out)

		fmt.Fprintf(w, "OK")
	}
}
