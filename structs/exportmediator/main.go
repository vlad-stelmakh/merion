package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// 'KZ241202-913860', 'KZ241203-684753', 'RU241201-567781', 'RU241201-795157'

func main() {
	ordersList := `
	'RU250121-746208', 'RU250121-563171', 'RU250121-578910', 'RU250121-975500', 'RU250121-949783', 'RU250121-770317', 'RU250121-770317'
	`

	ordersList = strings.TrimSpace(ordersList)
	ordersList = strings.ReplaceAll(ordersList, "'", "")
	orderIDs := strings.Split(ordersList, ",")

	for i := range orderIDs {
		orderIDs[i] = strings.TrimSpace(orderIDs[i])
	}

	url := "http://orders-export.services.lamoda.tech/jsonrpc/v2/outboxer.publish"

	payload := `{
            "id": "1",
            "jsonrpc": "2.0",
            "method": "outboxer.publish",
            "params": {
                "topic": "orders-export.order.updated",
                "payload": {
                    "order_nr": "%s"
                }
            }
        }`

	client := &http.Client{}

	for _, orderID := range orderIDs {
		formattedPayload := strings.NewReader(fmt.Sprintf(payload, orderID))

		req, err := http.NewRequest(http.MethodPost, url, formattedPayload)
		if err != nil {
			fmt.Println(err)
			return
		}

		req.Header.Add("Content-Type", "application/json")

		res, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Print("Procced order: " + orderID + " ")
		fmt.Println(string(body))

		time.Sleep(1 * time.Second)
	}
}
