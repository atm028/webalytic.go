package app

import (
	"encoding/json"
	"fmt"
	"net/http"

	Inversify "github.com/alekns/go-inversify"
	AppConfig "github.com/webalytic.go/cmd/collector/config"
	CommonCfg "github.com/webalytic.go/common/config"
	RedisBroker "github.com/webalytic.go/common/redis"
)

func CollectHandler(container Inversify.Container) http.HandlerFunc {
	cfgObj, _ := container.Get("appConfig")
	cfg, _ := cfgObj.(*AppConfig.CollectorConfig)
	cfgObj, _ = container.Get("redisBroker")
	broker, _ := cfgObj.(*RedisBroker.RedisBroker)
	logger := CommonCfg.GetLogger(container)

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
		logger.Debug("Endpoint collect payment: ", string(out))
		broker.Publish(cfg.Channel(), out)

		fmt.Fprintf(w, "OK")
	}
}
