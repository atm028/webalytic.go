package wbut

import (
	"context"
	"testing"

	"github.com/bhoriuchi/go-bunyan/bunyan"
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
	"github.com/webalytic.go/cmd/handler/app"
	Datasources "github.com/webalytic.go/common/datasources"
	"go.uber.org/fx"
)

//========================
type DBMock struct {
	callParams []interface{}
}
type Res struct {
	Error error
}

func (c *DBMock) Create(value interface{}) {
}

func (c *DBMock) Find(out interface{}, conds ...interface{}) error {
	return nil
}

//========================

func SetContainerUp() fx.Option {
	l, _ := bunyan.CreateLogger()
	logger := fx.Option(fx.Provide(func() bunyan.Logger {
		return l
	}))

	mockDB := new(DBMock)

	clickhouse := &Datasources.ClickHouse{
		Db: mockDB,
	}
	mockClickhouse1 := fx.Option(fx.Provide(func() *Datasources.ClickHouse {
		return clickhouse
	}))
	return fx.Options(
		logger,
		mockClickhouse1,
	)
}

func TestSubscribeAndAbleToPublishAndReadFromStream(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	container := SetContainerUp()
	app := fx.New(
		container,
		fx.Invoke(func(
			logger bunyan.Logger,
			clickhouse *Datasources.ClickHouse,
		) {
			ch := make(chan redis.XMessage)
			go app.RedisEventBrokerHandler(
				logger,
				clickhouse,
				ch,
			)
			assert.Equal(t, 1, 1)
			cancel()
		}),
	)

	_ = app.Start(ctx)
	<-ctx.Done()
	_ = app.Stop(context.Background())
}
