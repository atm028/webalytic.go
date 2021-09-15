package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/bhoriuchi/go-bunyan/bunyan"
	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/webalytic.go/cmd/collector/app"
	AppConfig "github.com/webalytic.go/cmd/collector/config"
	RedisBroker "github.com/webalytic.go/common/redis"
	"go.uber.org/fx"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	container := app.Container()
	mainApp := fx.New(
		container,
		fx.Invoke(func(
			appConfig *AppConfig.CollectorConfig,
			broker *RedisBroker.RedisBroker,
			httpCollectorHandler http.HandlerFunc,
			logger bunyan.Logger,
		) {
			collectorRedisChannel := make(chan redis.XMessage, 1)
			broker.Subscribe(appConfig.Channel(), collectorRedisChannel)

			router := mux.NewRouter().StrictSlash(true)
			router.HandleFunc("/collect", httpCollectorHandler).Methods("POST")

			//prometheus.MustRegister(prometheus.NewGoCollector())
			//prometheus.MustRegister(prometheus.NewBuildInfoCollector())

			router.Handle("/metrics", promhttp.HandlerFor(
				prometheus.DefaultGatherer,
				promhttp.HandlerOpts{
					EnableOpenMetrics: true,
				},
			))
			port := appConfig.Port()

			logger.Info("Started service on oprt ", port)
			log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
		}),
	)

	_ = mainApp.Start(ctx)
	<-ctx.Done()
	_ = mainApp.Stop(context.Background())
}
