package wbut

import (
	"context"
	"encoding/json"
	"fmt"
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

func TestRedisEventBroker(t *testing.T) {
	SetContainerUp := func() fx.Option {
		fmt.Println("==============SetContainerUp-------------")
		componentName := "test"
		commonCfgOption := CommonCfg.Container()
		commonOptions := Common.Container(componentName)
		redisBroker := RedisBroker.Container()
		return fx.Options(
			fx.Provide(func() string { return componentName }),
			commonOptions,
			commonCfgOption,
			redisBroker,
		)
	}
	/*
	   TODO: Need to understand why the onStop hook of EventBroker is not calling here.
	   Now it is a reason of the duplicate Prometheus metric registration
	*/
	//	t.Run("InitRedisBrokerCorrectly", func(t *testing.T) {
	//		container := SetContainerUp()
	//		ctx, cancel := context.WithCancel(context.Background())
	//		app := fx.New(
	//			container,
	//			fx.Invoke(func(
	//				logger bunyan.Logger,
	//				broker *RedisBroker.RedisBroker,
	//			) {
	//				assert.NotNil(t, broker)
	//				cancel()
	//			}),
	//		)
	//		_ = app.Start(ctx)
	//		<-ctx.Done()
	//		_ = app.Stop(context.Background())
	//	})

	t.Run("SubscribeAndAbleToPublishAndReadFromStream", func(t *testing.T) {
		container := SetContainerUp()
		ctx, cancel := context.WithCancel(context.Background())
		app := fx.New(
			container,
			fx.Invoke(func(
				logger bunyan.Logger,
				broker *RedisBroker.RedisBroker,
			) {
				ch := make(chan redis.XMessage)
				chName := "collector-stream"
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
	})

	t.Run("ReadFromLogHandlersAndAck", func(t *testing.T) {
		container := SetContainerUp()
		ctx, cancel := context.WithCancel(context.Background())
		app := fx.New(
			container,
			fx.Invoke(func(
				logger bunyan.Logger,
				broker *RedisBroker.RedisBroker,
			) {
				ch := make(chan redis.XMessage)
				chName := "service-stream"
				broker.Subscribe(chName, ch)
				msg := TestMsg{"1", "field1", "field2"}

				go func(ch chan redis.XMessage) {
					event := <-ch
					logger.Debug("received event from channel")
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
	})
}
