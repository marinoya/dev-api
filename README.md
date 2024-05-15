# Running

use ```go run .``` to run on localhost port 8088

# Constants and Variables

Constants such as MerchantID, APISecretKey, Currency, and URLs for the API endpoints define the basic configuration for interacting with the Zotapay API.
Variables like DepositRequests and OrderStatuses are used to store ongoing transactions and their statuses.

# HTTP Handlers

1. DepositHandler (POST /deposit): This handler initiates a deposit request to Zotapay. It reads the request, validates it, constructs a new request with necessary information, and sends it to Zotapay. If Zotapay accepts the request, it starts a background job to periodically check the status of the order until it reaches a final state.
2. CallbackHandler (POST /callback): Receives callback notifications from Zotapay regarding the status of a transaction. It verifies the signature of the incoming request to ensure it's genuine and updates the transaction status accordingly.
3. RedirectHandler (GET /payment_return): Handles redirection after a transaction is processed. It checks the signature of the incoming parameters, validates the transaction's existence, and then returns a message to the user based on the final status of the transaction (e.g., approved, declined).

# Utility Functions

checkStatus: Called periodically to check the current status of an order by making an API request to Zotapay.
GenerateSignature and ValidateSignature: Functions used to generate and validate signatures for secure communication with the Zotapay API. 