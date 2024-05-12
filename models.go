package main 

type DepositRequest struct {
    MerchantOrderID      string `json:"merchantOrderID"`
    MerchantOrderDesc    string `json:"merchantOrderDesc"`
    OrderAmount          string `json:"orderAmount"`
    OrderCurrency        string `json:"orderCurrency"`
    CustomerEmail        string `json:"customerEmail"`
    CustomerFirstName    string `json:"customerFirstName"`
    CustomerLastName     string `json:"customerLastName"`
    CustomerAddress      string `json:"customerAddress"`
    CustomerCountryCode  string `json:"customerCountryCode"`
    CustomerCity         string `json:"customerCity"`
    CustomerZipCode      string `json:"customerZipCode"`
    CustomerPhone        string `json:"customerPhone"`
    CustomerIP           string `json:"customerIP"`
    RedirectUrl          string `json:"redirectUrl"`
    CallbackUrl          string `json:"callbackUrl"`
    CustomParam          string `json:"customParam"`
    CheckoutUrl          string `json:"checkoutUrl"`
    Signature            string `json:"signature"`
}