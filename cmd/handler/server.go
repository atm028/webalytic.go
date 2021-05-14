package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/webalytic.go/cmd/handler/app"
	AppConfig "github.com/webalytic.go/cmd/handler/config"
	CommonCfg "github.com/webalytic.go/common/config"
	RedisBroker "github.com/webalytic.go/common/redis"
)

func main() {
	container := app.Container()
	cfgObj, _ := container.Get("appConfig")
	cfg, _ := cfgObj.(*AppConfig.LogHandlerConfig)
	cfgObj, _ = container.Get("redisBroker")
	broker, _ := cfgObj.(*RedisBroker.RedisBroker)

	collectorRedisChannel := make(chan redis.XMessage)
	broker.Subscribe(cfg.Channel(), collectorRedisChannel)
	go app.RedisEventBrokerHandler(container, collectorRedisChannel)

	router := mux.NewRouter().StrictSlash(true)
	port := cfg.Port()

	logger := CommonCfg.GetLogger(container)
	logger.Info("Started service at port: ", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
}
