package app

import (
	"fmt"
	"net/http"

	"github.com/bhoriuchi/go-bunyan/bunyan"
)

type IHealthHandler struct {
	Handler http.HandlerFunc
}

func HealthHandler(logger bunyan.Logger) *IHealthHandler {
	return &IHealthHandler{
		Handler: func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "passing")
		},
	}
}
