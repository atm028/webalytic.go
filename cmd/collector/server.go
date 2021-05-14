package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/webalytic.go/cmd/collector/app"
	AppConfig "github.com/webalytic.go/cmd/collector/config"
	CommonCfg "github.com/webalytic.go/common/config"
	RedisBroker "github.com/webalytic.go/common/redis"
)

func main() {
	container := app.Container()
	cfgObj, _ := container.Get("appConfig")
	cfg, _ := cfgObj.(*AppConfig.CollectorConfig)
	cfgObj, _ = container.Get("redisBroker")
	broker, _ := cfgObj.(*RedisBroker.RedisBroker)

	collectorRedisChannel := make(chan redis.XMessage, 1)
	broker.Subscribe(cfg.Channel(), collectorRedisChannel)

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/collect", app.CollectHandler(container)).Methods("POST")
	port := cfg.Port()

	logger := CommonCfg.GetLogger(container)
	logger.Info("Started service on oprt ", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
}
