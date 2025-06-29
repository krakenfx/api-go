// Prints the derivatives margin accounts.
package main

import (
	"fmt"
	"os"

	"github.com/krakenfx/api-go/v2/internal/helper"
	"github.com/krakenfx/api-go/v2/pkg/derivatives"
)

func main() {
	client := derivatives.NewREST()
	client.BaseURL = os.Getenv("KRAKEN_API_FUTURES_REST_URL")
	client.PublicKey = os.Getenv("KRAKEN_API_FUTURES_PUBLIC")
	client.PrivateKey = os.Getenv("KRAKEN_API_FUTURES_SECRET")
	fmt.Printf("> Fetching margin accounts.\n")
	accounts, err := client.Accounts()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Accounts: %s\n", helper.ToJSONIndent(accounts))
}
