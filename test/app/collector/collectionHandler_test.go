package wbut

import (
	"bytes"
	"context"
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/bhoriuchi/go-bunyan/bunyan"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/webalytic.go/cmd/collector/app"
	"go.uber.org/fx"
)

func TestHandlePaymentRequest(t *testing.T) {
	container := app.Container()
	ctx, cancel := context.WithCancel(context.Background())
	app := fx.New(
		container,
		fx.Invoke(func(
			logger bunyan.Logger,
			httpCollectorHandler *app.ICollectHandler,
		) {
			body := []byte(`
			{
				"Merchant": "macdonlds",
				"Sum": 169.89,
				"Project": "shop",
				"SendCurrency": "ru",
				"Method": "get",
				"Name": "some name",
				"CardName": "visa",
				"CardNumber": "12555667788",
				"ExpireDate": "11/22",
				"SecurityCode": "123",
				"ReceiveCurrency": "usd",
				"Rate": 65,
				"TransactionStatus": "done",
				"TransactionTime": "2012-04-23T18:25:43Z"
			}
			`)
			req, err := http.NewRequest("POST", "/collect/payment", bytes.NewReader(body))
			if err != nil {
				logger.Fatal(err)
			}
			rr := httptest.NewRecorder()
			router := mux.NewRouter()
			router.HandleFunc("/collect/{category}", httpCollectorHandler.Handler)
			router.ServeHTTP(rr, req)
			assert.Equal(t, rr.Code, http.StatusOK)
			cancel()
		}),
	)

	_ = app.Start(ctx)
	<-ctx.Done()
	_ = app.Stop(context.Background())
}
