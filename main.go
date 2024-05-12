package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

const MerchantID string = "BUGBOUNTY231"
const APISecretKey string = "5f4a6fcf-9048-4a0b-afc2-ed92d60fb1bf"
const Currency string = "USD"
const EndpointId string = "402334"
const BaseURL string = "https://api.zotapay-stage.com"

var DepositRequests = map[string]string{}
var OrderStatuses = map[string]string{}

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/deposit", DepositHandler).Methods("POST")

	router.HandleFunc("/status", StatusCheckHandler).Methods("GET")

	router.HandleFunc("/callback", CallbackHandler).Methods("POST")

	router.HandleFunc("/payment_return", RedirectHandler).Methods("GET")

	fmt.Println("Starting server on http://localhost:8088")

	if err := http.ListenAndServe(":8088", router); err != nil {
		fmt.Println(err)
	}

}

type ApiResponse struct {
	Status string `json:"status"`
}

// statusCheckFlow sends a GET request and returns the status from the response.
func statusCheckFlow() (string, error) {
	resp, err := http.Get("https://api.zotapay.com/api/v1/query/order-status/?${qs}")
	if err != nil {
		return "", err // Return an empty string and the error
	}
	defer resp.Body.Close()

	var apiResponse ApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return "", err // Return an empty string and the error
	}

	return apiResponse.Status, nil
}

func DepositHandler(w http.ResponseWriter, r *http.Request) {

	depositData := DepositRequest{
		MerchantOrderID:     "18",
		MerchantOrderDesc:   "Test order",
		OrderAmount:         "500.00",
		OrderCurrency:       Currency,
		CustomerEmail:       "customer@email-address.com",
		CustomerFirstName:   "John",
		CustomerLastName:    "Doe",
		CustomerAddress:     "5/5 Moo 5 Thong Nai Pan Noi Beach, Baan Tai, Koh Phangan",
		CustomerCountryCode: "TH",
		CustomerCity:        "Surat Thani",
		CustomerZipCode:     "84280",
		CustomerPhone:       "+66-77999110",
		CustomerIP:          "103.106.8.104",
		RedirectUrl:         "https://www.example-merchant.com/payment-return/",
		CallbackUrl:         "https://www.example-merchant.com/payment-callback/",
		CustomParam:         "{\"UserId\": \"e139b447\"}",
		CheckoutUrl:         "https://www.example-merchant.com/account/deposit/?uid=e139b447",
		Signature:           GenerateSignature(EndpointId + "18" + "500" + "customer@email-address.com" + APISecretKey),
	}

	jsonData, err := json.Marshal(depositData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	request, err := http.NewRequest("POST", BaseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	fmt.Fprintf(w, "Deposit flow initiated")
}

func StatusCheckHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "Status check flow initiated")
}

func CallbackHandler(w http.ResponseWriter, r *http.Request) {

	var requestParams CallbackModel

	if body, err := io.ReadAll(r.Body); err != nil {
		fmt.Println("Error reading body", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		if err = json.Unmarshal(body, &requestParams); err != nil {
			fmt.Println("Error parsing JSON", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	params := requestParams.EndpointID + requestParams.OrderID + requestParams.Status + requestParams.Amount + requestParams.CustomerEmail + APISecretKey
	if !ValidateSignature(params, requestParams.Signature) {
		w.WriteHeader(http.StatusUnauthorized)

		return
	}

	OrderStatuses[requestParams.OrderID] = requestParams.Status

	w.WriteHeader(http.StatusOK)
}

func RedirectHandler(w http.ResponseWriter, r *http.Request) {

	q := r.URL.Query()
	status := q.Get("status")
	orderID := q.Get("orderID")
	merchantOrderID := q.Get("merchantOrderID")
	signature := q.Get("signature")

	if !ValidateSignature(status+orderID+merchantOrderID+APISecretKey, signature) {
		fmt.Printf("Invalid signature for parameters: `%s`, `%s`, `%s`, provided signature was `%s`", status, orderID, merchantOrderID, signature)
		w.WriteHeader(http.StatusUnauthorized)

		return
	}

	oid, exists := DepositRequests[merchantOrderID]

	if !exists {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	orderStatus := OrderStatuses[oid]
	var result string

	switch orderStatus {
	case "APPROVED":
		result = "Success, thank you for your painment!"
	case "DECLINED", "FILTERED", "ERROR":
		result = "Failed, please try again!"
	default:
		result = "Please wait"
	}

	w.Write([]byte(result))
	w.WriteHeader(http.StatusOK)

	fmt.Fprintf(w, "Status check flow initiated")
}
