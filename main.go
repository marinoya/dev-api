package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

const MerchantID = "BUGBOUNTY231"
const APISecretKey = "5f4a6fcf-9048-4a0b-afc2-ed92d60fb1bf"
const Currency = "USD"
const EndpointId = "402334"
const BaseURL = "https://api.zotapay-stage.com"
const ZotaOrderAcceptedCode = "200"

var ZotaOrderFinalStates = []string{"APPROVED", "DECLINED", "FILTERED", "ERROR"}

const pollingTime = 10 * time.Second
const depositPath = "/api/v1/deposit/request/" + EndpointId + "/"
const orderStatusPath = "/api/v1/query/order-status/"

const ServerURL = "http://localhost"
const ServerPort = ":8088"
const CallbackPath = "/callback"
const PaymentPath = "/payment_return"
const DepositPath = "/deposit"

var DepositRequests = map[string]string{}
var OrderStatuses = map[string]string{}

func main() {

	router := mux.NewRouter()

	router.HandleFunc(DepositPath, DepositHandler).Methods("POST")

	router.HandleFunc(CallbackPath, CallbackHandler).Methods("POST")

	router.HandleFunc(PaymentPath, RedirectHandler).Methods("GET")

	fmt.Println("Starting server on " + ServerURL + ServerPort)

	if err := http.ListenAndServe(ServerPort, router); err != nil {
		fmt.Println(err)
	}

}

func DepositHandler(w http.ResponseWriter, r *http.Request) {

	var requestParams DepositModel

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

	if requestParams.OrderCurrency != Currency {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	depositData := DepositRequest{
		MerchantOrderID:     uuid.New().String(),
		MerchantOrderDesc:   requestParams.MerchantOrderDesc,
		OrderAmount:         requestParams.OrderAmount,
		OrderCurrency:       Currency,
		CustomerEmail:       requestParams.CustomerEmail,
		CustomerFirstName:   requestParams.CustomerFirstName,
		CustomerLastName:    requestParams.CustomerLastName,
		CustomerAddress:     requestParams.CustomerAddress,
		CustomerCountryCode: requestParams.CustomerCountryCode,
		CustomerCity:        requestParams.CustomerCity,
		CustomerZipCode:     requestParams.CustomerZipCode,
		CustomerPhone:       requestParams.CustomerPhone,
		CustomerIP:          r.RemoteAddr,
		RedirectUrl:         ServerURL + ServerPort + PaymentPath,
		CallbackUrl:         ServerURL + ServerPort + CallbackPath,
		CheckoutUrl:         ServerURL + ServerPort + DepositPath,
	}

	depositData.Signature = GenerateSignature(EndpointId + depositData.MerchantOrderID + depositData.OrderAmount + depositData.CustomerEmail + APISecretKey)

	jsonData, err := json.Marshal(depositData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	request, err := http.NewRequest("POST", BaseURL+depositPath, bytes.NewBuffer(jsonData))
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

	var depositResponse DepositResponse

	if body, err := io.ReadAll(response.Body); err != nil {
		fmt.Println("Error reading response body", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		if err = json.Unmarshal(body, &depositResponse); err != nil {
			fmt.Println("Error parsing response JSON", body, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	if depositResponse.Code != ZotaOrderAcceptedCode {
		fmt.Println("Zota Error", depositResponse.Code, depositResponse.Message)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	DepositRequests[depositData.MerchantOrderID] = depositResponse.Data.OrderID

	go func() {

		for {

			time.Sleep(pollingTime)

			currentStatus := OrderStatuses[depositResponse.Data.OrderID]

			if slices.Contains(ZotaOrderFinalStates, currentStatus) {
				return
			}

			checkStatus(depositData.MerchantOrderID, depositResponse.Data.OrderID)
		}
	}()

	w.Write([]byte(depositResponse.Data.DepositUrl))
	w.WriteHeader(http.StatusFound)

}

func checkStatus(moid, oid string) {
	t := fmt.Sprint(time.Now().Unix())
	s := GenerateSignature(MerchantID + moid + oid + t + APISecretKey)

	requestURL := fmt.Sprintf("%s%s?merchantID=%s&merchantOrderID=%s&orderID=%s&timestamp=%s&signature=%s", BaseURL, orderStatusPath, MerchantID, moid, oid, t, s)

	response, _ := http.Get(requestURL)

	var statusResponse CheckStatusResponse

	if body, err := io.ReadAll(response.Body); err == nil {
		if err = json.Unmarshal(body, &statusResponse); err == nil && statusResponse.Code == ZotaOrderAcceptedCode {
			OrderStatuses[oid] = statusResponse.Data.Status
		}
	}
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

	if slices.Contains(ZotaOrderFinalStates, status) {
		OrderStatuses[oid] = status
	}

	orderStatus := OrderStatuses[oid]
	var result string

	switch orderStatus {
	case "APPROVED":
		result = "Success, thank you for your paiment!"
	case "DECLINED", "FILTERED", "ERROR":
		result = "Failed, please try again!"
	default:
		result = "Please wait"
	}

	w.Write([]byte(result))
	w.WriteHeader(http.StatusOK)
}
