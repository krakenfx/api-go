package main

import (
	"fmt"
	"os"

	"github.com/krakenfx/api-go/pkg/spot"
)

func main() {
	client := spot.NewREST()
	client.BaseURL = os.Getenv("KRAKEN_API_SPOT_REST_URL")
	assetManager := spot.NewAssetManager()
	if err := assetManager.Use(client); err != nil {
		panic(err)
	}
	for _, alias := range []string{"XXBTZUSD", "XBTUSD", "XXBT/ZUSD", "XXBT", "XBT", "btc", "XDG", "XDG/USD"} {
		fmt.Printf("%s -> %s\n", alias, assetManager.Name(alias))
	}
}
