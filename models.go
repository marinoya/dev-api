package main

type DepositModel struct {
	MerchantOrderDesc   string
	OrderAmount         string
	OrderCurrency       string
	CustomerEmail       string
	CustomerFirstName   string
	CustomerLastName    string
	CustomerAddress     string
	CustomerCountryCode string
	CustomerCity        string
	CustomerZipCode     string
	CustomerPhone       string
}

type DepositRequest struct {
	MerchantOrderID     string
	MerchantOrderDesc   string
	OrderAmount         string
	OrderCurrency       string
	CustomerEmail       string
	CustomerFirstName   string
	CustomerLastName    string
	CustomerAddress     string
	CustomerCountryCode string
	CustomerCity        string
	CustomerZipCode     string
	CustomerPhone       string
	CustomerIP          string
	RedirectUrl         string
	CallbackUrl         string
	CheckoutUrl         string
	Signature           string
}

type DepositResponse struct {
	Code    string
	Message *string
	Data    *DataResponse
}

type DataResponse struct {
	DepositUrl      string
	MerchantOrderID string
	OrderID         string
}

type CallbackModel struct {
	EndpointID      string
	OrderID         string
	MerchantOrderID string
	Status          string
	Amount          string
	CustomerEmail   string
	Signature       string
}
