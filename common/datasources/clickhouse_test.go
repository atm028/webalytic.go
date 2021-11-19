package Datasources

import (
	"context"
	"testing"
	"time"

	"github.com/bhoriuchi/go-bunyan/bunyan"
	"github.com/stretchr/testify/assert"

	Common "github.com/webalytic.go/common"
	CommonCfg "github.com/webalytic.go/common/config"
	"go.uber.org/fx"
)

func TestClickhouse(t *testing.T) {
	Container := func() fx.Option {
		componentName := "test"
		commonCfgOptions := CommonCfg.Container()
		commonOptions := Common.Container(componentName)
		clickhouse := Clickhouse()
		ackRedisChannel := fx.Option(fx.Provide(func() chan string {
			ch := make(chan string)
			return ch
		}))

		return fx.Options(
			commonCfgOptions,
			commonOptions,
			clickhouse,
			ackRedisChannel,
		)
	}

	t.Run("CreateTableAndInsertRecord", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		container := Container()
		app := fx.New(
			container,
			fx.Invoke(func(
				logger bunyan.Logger,
				clickhouse *ClickHouse,
			) {
				payment_in := &Payment{
					TraceID:           "1",
					Merchant:          "TestMerchant",
					Sum:               123.45,
					SendCurrency:      "US",
					Project:           "TestProject",
					Method:            "card",
					Name:              "TestName",
					CardNumber:        "1234567890",
					ExpireDate:        "2050-01-01",
					SecurityCode:      "123",
					ReceiveCurrency:   "USD",
					Rate:              12,
					TransactionTime:   time.Now(),
					TransactionStatus: "in_progress",
					Field1:            "test_field1",
					Field2:            "test_field2",
					Field3:            "test_field3",
					Field4:            "test_field4",
					Field5:            "tests_field5",
					Field6:            "tests_field6",
					Field7:            "tests_field7",
					Field8:            "tests_field8",
					Field9:            "tests_field9",
					Field10:           "tests_field10",
				}

				clickhouse.CreatePayment("key1", payment_in)
				time.Sleep(time.Second)
				payment_out, _ := clickhouse.FindPayment("trace_id=?", "1")
				assert.Equal(t, payment_in.Sum, payment_out.Sum)
				assert.Equal(t, payment_in.SendCurrency, payment_out.SendCurrency)
				assert.Equal(t, payment_in.Project, payment_out.Project)
				assert.Equal(t, payment_in.Method, payment_out.Method)
				assert.Equal(t, payment_in.Name, payment_out.Name)
				assert.Equal(t, payment_in.CardNumber, payment_out.CardNumber)
				assert.Equal(t, payment_in.ExpireDate, payment_out.ExpireDate)
				assert.Equal(t, payment_in.SecurityCode, payment_out.SecurityCode)
				assert.Equal(t, payment_in.ReceiveCurrency, payment_out.ReceiveCurrency)
				assert.Equal(t, payment_in.Rate, payment_out.Rate)
				assert.Equal(t, payment_in.TransactionStatus, payment_out.TransactionStatus)
				assert.Equal(t, payment_in.Field1, payment_out.Field1)
				assert.Equal(t, payment_in.Field2, payment_out.Field2)
				assert.Equal(t, payment_in.Field3, payment_out.Field3)
				assert.Equal(t, payment_in.Field4, payment_out.Field4)
				assert.Equal(t, payment_in.Field5, payment_out.Field5)
				assert.Equal(t, payment_in.Field6, payment_out.Field6)
				assert.Equal(t, payment_in.Field7, payment_out.Field7)
				assert.Equal(t, payment_in.Field8, payment_out.Field8)
				assert.Equal(t, payment_in.Field9, payment_out.Field9)
				assert.Equal(t, payment_in.Field10, payment_out.Field10)
				clickhouse.Db.Migrator().DropTable(&Payment{})
				cancel()
			}),
		)
		_ = app.Start(ctx)
		<-ctx.Done()
		_ = app.Stop(context.Background())
	})
}
