package wbut

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/bhoriuchi/go-bunyan/bunyan"
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
	Common "github.com/webalytic.go/common"

	CommonCfg "github.com/webalytic.go/common/config"
	RedisBroker "github.com/webalytic.go/common/redis"
	"go.uber.org/fx"
)

type TestMsg struct {
	ID     string `json: "ID"`
	Field1 string `json: "Field1`
	Field2 string `json: "Field2`
}

func SetContainerUp() fx.Option {
	componentName := "test"
	commonCfgOption := CommonCfg.Container()
	commonOptions := Common.Container(componentName)
	redisBroker := RedisBroker.Container()
	return fx.Options(
		commonOptions,
		commonCfgOption,
		redisBroker,
	)
}

func TestInitRedisBrokerCorrectly(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	container := SetContainerUp()
	app := fx.New(
		container,
		fx.Invoke(func(
			logger bunyan.Logger,
			broker *RedisBroker.RedisBroker,
		) {
			assert.NotNil(t, broker)
			cancel()
		}),
	)
	_ = app.Start(ctx)
	<-ctx.Done()
	_ = app.Stop(context.Background())
}

func TestSubscribeAndAbleToPublishAndReadFromStream(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	container := SetContainerUp()
	app := fx.New(
		container,
		fx.Invoke(func(
			logger bunyan.Logger,
			broker *RedisBroker.RedisBroker,
		) {
			ch := make(chan redis.XMessage)
			chName := "testChannel"
			broker.Subscribe(chName, ch)
			msg := TestMsg{"1", "fild1", "field2"}

			go func(ch chan redis.XMessage, msg TestMsg) {
				logger.Info("Test:channel handler: start")
				event := <-ch
				logger.Info("Test:channel handler: get event")
				v, exist := event.Values["msg"]
				assert.Equal(t, exist, true)
				out, ok := v.(string)
				assert.Equal(t, ok, true)
				var rcv TestMsg
				err := json.Unmarshal([]byte(out), &rcv)
				assert.Nil(t, err)
				assert.Equal(t, msg, rcv)
				cancel()
			}(ch, msg)

			out, err := json.Marshal(msg)
			assert.Nil(t, err)
			broker.Publish(chName, out)
		}),
	)

	_ = app.Start(ctx)
	<-ctx.Done()
	_ = app.Stop(context.Background())
}

func TestReadFromLogHandlersAndAck(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	container := SetContainerUp()
	app := fx.New(
		container,
		fx.Invoke(func(
			logger bunyan.Logger,
			broker *RedisBroker.RedisBroker,
		) {
			ch := make(chan redis.XMessage)
			chName := "test"
			broker.Subscribe(chName, ch)
			msg := TestMsg{"1", "field1", "field2"}

			go func(ch chan redis.XMessage) {
				event := <-ch
				v, exist := event.Values["msg"]
				assert.Equal(t, exist, true)
				out, ok := v.(string)
				assert.Equal(t, ok, true)
				var rcv TestMsg
				err := json.Unmarshal([]byte(out), &rcv)
				assert.Nil(t, err)
				assert.Equal(t, msg, rcv)
				broker.Ack(chName, event.ID)
				cancel()
			}(ch)

			out, err := json.Marshal(msg)
			assert.Nil(t, err)
			err = broker.Publish(chName, out)
			assert.Nil(t, err)
		}),
	)
	_ = app.Start(ctx)
	<-ctx.Done()
	_ = app.Stop(context.Background())
}
