package app

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bhoriuchi/go-bunyan/bunyan"
	AppConfig "github.com/webalytic.go/cmd/collector/config"
	CommonCfg "github.com/webalytic.go/common/config"
	RedisBroker "github.com/webalytic.go/common/redis"
)

func CollectHandler(
	logger bunyan.Logger, broker *RedisBroker.RedisBroker,
	cfg *AppConfig.CollectorConfig,
	redisCfg *CommonCfg.RedisConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Debug("Collect handler")
		var payment Payment
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
		out, err := json.Marshal(payment)
		if err != nil {
			logger.Error(err)
		}
		logger.Debug("Endpoint collect payment: %s, channel: %s", string(out), redisCfg.StreamName())
		broker.Publish(redisCfg.StreamName(), out)

		fmt.Fprintf(w, "OK")
	}
}
