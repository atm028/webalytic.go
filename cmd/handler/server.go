package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/bhoriuchi/go-bunyan/bunyan"
	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/webalytic.go/cmd/handler/app"
	AppConfig "github.com/webalytic.go/cmd/handler/config"
	CommonCfg "github.com/webalytic.go/common/config"
	RedisBroker "github.com/webalytic.go/common/redis"
	"go.uber.org/fx"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	container := app.Container()
	mainApp := fx.New(
		container,
		fx.Invoke(func(
			appConfig *AppConfig.LogHandlerConfig,
			redisCfg *CommonCfg.RedisConfig,
			broker *RedisBroker.RedisBroker,
			logger bunyan.Logger,
		) {
			collectorRedisChannel := make(chan redis.XMessage)
			broker.Subscribe(redisCfg.StreamName(), collectorRedisChannel)
			go app.RedisEventBrokerHandler(logger, broker, appConfig, collectorRedisChannel)

			router := mux.NewRouter().StrictSlash(true)
			port := appConfig.Port()

			logger.Info("Started service at port: ", port)
			log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
		}),
	)

	_ = mainApp.Start(ctx)
	<-ctx.Done()
	_ = mainApp.Stop(context.Background())
}
