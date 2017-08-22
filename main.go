package main

import (
	"fmt"
	"time"
	"net/http"
	"log"
	"encoding/json"
	"github.com/cactus/go-statsd-client/statsd"
	"strconv"
	"errors"
	"os"
)

type BillingList struct {
	Invoices []struct {
		TotalUntaxed      string      `json:"total_untaxed"`
		OrganizationName  string      `json:"organization_name"`
		DueDate           interface{} `json:"due_date"`
		TotalTax          string      `json:"total_tax"`
		StopDate          time.Time   `json:"stop_date"`
		IssuedDate        interface{} `json:"issued_date"`
		Number            interface{} `json:"number"`
		ID                string      `json:"id"`
		OrganizationID    string      `json:"organization_id"`
		Currency          string      `json:"currency"`
		State             string      `json:"state"`
		TotalTaxed        string      `json:"total_taxed"`
		LastUpdate        time.Time   `json:"last_update"`
		TotalUndiscounted string      `json:"total_undiscounted"`
		StartDate         time.Time   `json:"start_date"`
	} `json:"invoices"`
}

func main() {

	statsdHost := getenv("STATSD_ADDRESS", "127.0.0.1:8125")

	statsdPrefix := getenv("STATSD_PREFIX", "scaleway-billing")

	checkDelay, err := strconv.ParseInt(getenv("CHECK_DELAY_SECONDS", "10"), 10, 32)

	if err != nil {
		log.Fatal("Error parsing check delay parameter", err)
		return
	}

	authToken := getenv("SCALEWAY_AUTH_TOKEN", "")

	if authToken == "" {
		log.Fatal("Scaleway auth token is empty.")
		return
	}

	statsDClient, _ := statsd.NewClient(statsdHost, statsdPrefix)

	client := &http.Client{}

	defer statsDClient.Close()

	for ; ; {

		response, err := getBillingList(client, authToken)
		if err != nil {
			log.Fatal("Error getting billing list.")
			return
		}

		sendData(statsDClient, response)

		time.Sleep(time.Duration(checkDelay) * time.Second)
	}

}

func sendData(statsd statsd.Statter, list BillingList) {
	invoice := list.Invoices[0]

	fmt.Printf("Invoice ID: %s\n", invoice.ID)

	totalTaxed, err := fstringToInt64(invoice.TotalTaxed)
	totalTax, err := fstringToInt64(invoice.TotalTax)
	totalUntaxed, err := fstringToInt64(invoice.TotalUntaxed)
	totalUndiscounted, err := fstringToInt64(invoice.TotalUndiscounted)

	valuesMap := make(map[string]int64)
	valuesMap["current.total-taxed"] = totalTaxed
	valuesMap["current.total-tax"] = totalTax
	valuesMap["current.total-untaxed"] = totalUntaxed
	valuesMap["current.total-undiscounted"] = totalUndiscounted

	for k, v := range valuesMap {
		fmt.Println(k, ": ", v)
		if statsd.Gauge(k, v, 1.0) != nil {
			log.Println(err)
		}
	}

}

func getBillingList(client *http.Client, authToken string) (BillingList, error) {

	apiUrl := getBillingListUrl(authToken)

	req, _ := http.NewRequest("GET", apiUrl, nil)

	resp, err := client.Do(req)
	defer resp.Body.Close()

	var response BillingList

	if err != nil {
		log.Fatal("Error executing request to "+apiUrl, err)
		return response, errors.New("Error executing http request.")
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Println(err)
		return response, errors.New("Error parsing json.")
	}

	return response, nil
}

func fstringToInt64(in string) (int64, error) {
	floatValue, err := strconv.ParseFloat(in, 32)
	if err != nil {
		log.Fatal("Float parse error", err)
		return 0, errors.New("Invalid float string " + in)
	}
	return int64(floatValue * 100), nil

}

func getBillingListUrl(authToken string) string {
	return fmt.Sprintf("https://billing.scaleway.com/invoices/?x-auth-token=%s", authToken)
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
