package main

import (
	"fmt"
	"os"

	"github.com/krakenfx/api-go/pkg/derivatives"
	"github.com/krakenfx/api-go/pkg/kraken"
)

func main() {
	client := derivatives.NewREST()
	client.BaseURL = os.Getenv("KRAKEN_API_FUTURES_REST_URL")
	client.PublicKey = os.Getenv("KRAKEN_API_FUTURES_PUBLIC")
	client.PrivateKey = os.Getenv("KRAKEN_API_FUTURES_SECRET")
	fmt.Printf("> Fetching unwind queue.\n")
	response, err := client.Issue(&derivatives.RequestOptions{
		Method: "GET",
		Path:   []any{"/derivatives/api/v3/unwindqueue"},
		Auth:   true,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Unwind queue: %s\n", kraken.ToJSONIndent(response.BodyMap))
}
