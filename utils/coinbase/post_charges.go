package coinbase

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type CheckoutResponse struct {
	Data struct {
		ID string `json:"id"`
	} `json:"data"`
}

func PostCharges(userInput int) string {

	amount := userInput
	if amount == 0 {
		amount = 100
	}
	url := "https://api.commerce.coinbase.com/checkouts"
	method := "POST"
	name := "BUY DERO WITH CRYPTO"
	description := "Purchase DERO using Coinbase's Commerce Platform. \n Once we have delivered, you will have DERO \n You are being asked for your DERO wallet address/name"
	requestedInfo := []string{"address"}
	pricingType := "fixed_price"
	localPrice := map[string]string{"amount": fmt.Sprintf("%d", amount), "currency": "USD"}
	coinbaseAPIToken := os.Getenv("COINBASE_API_TOKEN")

	payloadData := map[string]interface{}{
		"name":           name,
		"description":    description,
		"requested_info": requestedInfo,
		"pricing_type":   pricingType,
		"local_price":    localPrice,
	}

	payloadJSON, err := json.Marshal(payloadData)
	if err != nil {
		fmt.Println(err)
		return "ERROR"
	}

	payload := strings.NewReader(string(payloadJSON))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		fmt.Println(err)
		return "ERROR"
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-CC-Version", "2018-03-22")
	req.Header.Add("X-CC-Api-Key", coinbaseAPIToken)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "ERROR"
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return "ERROR"
	}

	var checkoutResponse CheckoutResponse
	if err := json.Unmarshal(body, &checkoutResponse); err != nil {
		fmt.Println(err)
		return "ERROR"
	}
	Checkout := checkoutResponse.Data.ID
	fmt.Println("Checkout ID:", Checkout)
	return Checkout
}
