package Datasources

import "time"

type Payment struct {
	TraceID           string    `gorm:"primaryKey;autoIncrement:true;unique"`
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
	Field1            string    `json: "fied1"`
	Field2            string    `json: "fied2"`
	Field3            string    `json: "fied3"`
	Field4            string    `json: "fied4"`
	Field5            string    `json: "fied5"`
	Field6            string    `json: "fied6"`
	Field7            string    `json: "fied7"`
	Field8            string    `json: "fied8"`
	Field9            string    `json: "fied9"`
	Field10           string    `json: "fied10"`
}
