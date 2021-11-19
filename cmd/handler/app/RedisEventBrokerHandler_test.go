package app

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/bhoriuchi/go-bunyan/bunyan"
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
	Datasources "github.com/webalytic.go/common/datasources"
	"go.uber.org/fx"
)

func SetContainerUp() fx.Option {
	container := Container()
	return container
}

func TestSubscribeAndAbleToPublishAndReadFromStream(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	container := SetContainerUp()
	traceID := fmt.Sprintf("%d", rand.Intn(10000))
	payment := Datasources.Payment{
		TraceID:           traceID,
		Merchant:          "kamwnlsqjf",
		Sum:               60,
		SendCurrency:      "ru",
		Project:           "qfdjstnulw",
		Method:            "get",
		Name:              "isdcvmcusz",
		CardNumber:        "920236688367",
		ExpireDate:        "11/22",
		SecurityCode:      "123",
		ReceiveCurrency:   "usd",
		Rate:              65,
		TransactionTime:   time.Now(),
		TransactionStatus: "done",
		Field1:            "",
		Field2:            "",
		Field3:            "",
		Field4:            "",
		Field5:            "",
		Field6:            "",
		Field7:            "",
		Field8:            "",
		Field9:            "",
		Field10:           "",
	}

	app := fx.New(
		container,
		fx.Invoke(func(
			logger bunyan.Logger,
			clickhouse *Datasources.ClickHouse,
		) {
			out, _ := json.Marshal(&payment)
			redisMsg := redis.XMessage{
				ID: "testMessage1",
				Values: map[string]interface{}{
					"msg": string(out[:]),
				},
			}

			ch := make(chan redis.XMessage)
			go RedisEventBrokerHandler(logger, clickhouse, ch)
			ch <- redisMsg
			time.Sleep(2 * time.Second)
			foundPayment, err := clickhouse.FindPayment("trace_id", traceID)
			assert.Equal(t, err, nil)
			assert.Equal(t, payment.TraceID, foundPayment.TraceID)
			cancel()
		}),
	)

	_ = app.Start(ctx)
	<-ctx.Done()
	_ = app.Stop(context.Background())
}
