package main

import (
	"fmt"
	"os"

	"github.com/krakenfx/api-go/v2/internal/helper"
	"github.com/krakenfx/api-go/v2/pkg/spot"
)

func main() {
	client := spot.NewREST()
	client.BaseURL = os.Getenv("KRAKEN_API_SPOT_REST_URL")
	client.PublicKey = os.Getenv("KRAKEN_API_SPOT_PUBLIC")
	client.PrivateKey = os.Getenv("KRAKEN_API_SPOT_SECRET")
	withdrawMethods, err := client.Call(spot.RequestOptions{
		Auth:   true,
		Method: "POST",
		Path:   []any{"/0/private/WithdrawMethods"},
		Body: map[string]any{
			"asset": "XRP",
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Withdraw methods: %s\n", helper.ToJSONIndent(withdrawMethods))
}
