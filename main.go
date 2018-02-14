package main

import (
	"time"
	"os"
	"github.com/GaruGaru/scaleway-billing-monitor/api"
	"github.com/GaruGaru/scaleway-billing-monitor/stats"
	"fmt"
)

func main() {

	authToken, present := os.LookupEnv("SCALEWAY_AUTH_TOKEN")

	if !present {
		panic("Empty env. var SCALEWAY_AUTH_TOKEN")
	}

	billingApi := api.BillingApi{Token: authToken}

	billingStatter := stats.StatterHolderOf(stats.NewStatsdStatterWithFallback())

	for ; ; {

		billing, err := billingApi.BillingList()

		if err != nil {
			fmt.Printf("Unable to fetch billing info: %s\n", err.Error())
		} else {
			fmt.Printf("Got %d results from api\n", len(billing.Invoices))
			billingStatter.Send(billing)
		}

		time.Sleep(10 * time.Second)
	}

}
