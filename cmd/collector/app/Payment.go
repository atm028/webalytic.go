package app

import "time"

type Payment struct {
	Merchant          string    `json: "Merchant"`
	Sum               float64   `json: "Sum"`
	SendCurrency      string    `json: "SendCurrency"`
	Project           string    `json: "Project"`
	Method            string    `json: "Method"`
	Name              string    `json: "Name"`
	CardNumber        string    `json: "CardNumber"`
	ExpireDate        string    `json: "ExpireDate"`
	SecurityCode      string    `json: "SecurityCode"`
	ReceiveCurrency   string    `json: "ReceiveCurrency"`
	Rate              float64   `json: "Rate"`
	TransactionTime   time.Time `json: "TransactionTime"`
	TransactionStatus string    `json: "TransactionStatus"`
	field1            string    `json: "fied1"`
	field2            string    `json: "fied2"`
	field3            string    `json: "fied3"`
	field4            string    `json: "fied4"`
	field5            string    `json: "fied5"`
	field6            string    `json: "fied6"`
	field7            string    `json: "fied7"`
	field8            string    `json: "fied8"`
	field9            string    `json: "fied9"`
	field10           string    `json: "fied10"`
}
