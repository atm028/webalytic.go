package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/bhoriuchi/go-bunyan/bunyan"
	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/webalytic.go/cmd/handler/app"
	AppConfig "github.com/webalytic.go/cmd/handler/config"
	CommonCfg "github.com/webalytic.go/common/config"
	Datasources "github.com/webalytic.go/common/datasources"
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
			appConfig *AppConfig.HandlerConfig,
			redisCfg *CommonCfg.RedisConfig,
			broker *RedisBroker.RedisBroker,
			clickhouse *Datasources.ClickHouse,
			logger bunyan.Logger,
			ackRedisChannel chan string,
		) {
			collectorRedisChannel := make(chan redis.XMessage)
			broker.Subscribe(redisCfg.StreamName(), collectorRedisChannel)
			go app.RedisEventBrokerHandler(
				logger,
				broker,
				clickhouse,
				appConfig,
				collectorRedisChannel)

			//Channel for acknowledgement handled messages from redis when they are saved into the DB
			go func() {
				for {
					key := <-ackRedisChannel
					if len(key) != 0 {
						broker.GroupAck(key)
					}
				}
			}()

			router := mux.NewRouter().StrictSlash(true)

			router.Handle("/metrics", promhttp.HandlerFor(
				prometheus.DefaultGatherer,
				promhttp.HandlerOpts{
					EnableOpenMetrics: true,
				},
			))
			port := appConfig.Port()

			logger.Info(fmt.Sprintf("Started service at port: %d", port))
			log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
		}),
	)

	_ = mainApp.Start(ctx)
	<-ctx.Done()
	_ = mainApp.Stop(context.Background())
}
