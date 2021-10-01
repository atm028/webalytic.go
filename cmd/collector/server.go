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
	Consul "github.com/webalytic.go/common/consul"
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
			httpCollectorHandler *app.ICollectHandler,
			httpHealthHandler *app.IHealthHandler,
			logger bunyan.Logger,
			consul *Consul.Consul,
			name string,
		) {
			collectorRedisChannel := make(chan redis.XMessage, 1)
			broker.Subscribe(appConfig.Channel(), collectorRedisChannel)

			router := mux.NewRouter().StrictSlash(true)
			router.HandleFunc("/collect/{category}", httpCollectorHandler.Handler).Methods("POST")

			router.Handle("/metrics", promhttp.HandlerFor(
				prometheus.DefaultGatherer,
				promhttp.HandlerOpts{
					EnableOpenMetrics: true,
				},
			))

			port := appConfig.Port()
			router.Handle("/info", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				logger.Debug(fmt.Sprintf("Info request to %s", name))
				fmt.Fprintf(w, name)
			})).Methods("GET")
			router.Handle("/health", httpHealthHandler.Handler).Methods("GET")
			consul.ServiceRegister(name, port, "5s", "/health", "GET")

			logger.Info(fmt.Sprintf("Started service on port: %d", port))
			log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
		}),
	)

	_ = mainApp.Start(ctx)
	<-ctx.Done()
	_ = mainApp.Stop(context.Background())
}
