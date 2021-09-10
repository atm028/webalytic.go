package app

import (
	"encoding/json"
	"fmt"

	"github.com/bhoriuchi/go-bunyan/bunyan"
	"github.com/go-redis/redis"
	AppConfig "github.com/webalytic.go/cmd/handler/config"
	Datasources "github.com/webalytic.go/common/datasources"
	RedisBroker "github.com/webalytic.go/common/redis"
)

func RedisEventBrokerHandler(
	logger bunyan.Logger,
	broker *RedisBroker.RedisBroker,
	clickhouse *Datasources.ClickHouse,
	cfg *AppConfig.LogHandlerConfig,
	evtChannel chan redis.XMessage) {
	for {
		event := <-evtChannel
		logger.Debug(event.Values)
		v, exist := event.Values["msg"]
		if !exist {
			logger.Error(fmt.Sprintf("Error: No message found for event ID: %s", event.ID))
		} else {
			out, ok := v.(string)
			if !ok {
				logger.Error(fmt.Sprintf("Unable to convert message to string for event with ID: %s", event.ID))
			} else {
				var rcv Datasources.Payment
				err := json.Unmarshal([]byte(out), &rcv)
				if err != nil {
					logger.Error(fmt.Sprintf("Unable to parse to JSON message from event with ID: %s", event.ID))
				} else {
					logger.Debug(fmt.Sprintf("traceID: %s: ID: %s: %v", rcv.TraceID, event.ID, rcv))
					rcv.CacheKey = event.ID
					clickhouse.CreatePayment(&rcv)
				}
			}
		}
	}
}
