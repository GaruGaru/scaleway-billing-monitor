package api

import (
	"fmt"
	"time"
	"net/http"
	"encoding/json"
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

type BillingApi struct {
	Token string
}

func (api BillingApi) getBillingListUrl() string {
	return fmt.Sprintf("https://billing.scaleway.com/invoices/?x-auth-token=%s", api.Token)
}

func (api BillingApi) BillingList() (BillingList, error) {
	resp, err := http.Get(api.getBillingListUrl())

	if err != nil {
		return BillingList{}, err
	}

	defer resp.Body.Close()

	var response BillingList

	if resp.StatusCode != 200 {
		return BillingList{}, fmt.Errorf("scaleway api status code error: %d\n", resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(&response)

	if err != nil {
		return BillingList{}, err
	}

	return response, nil

}
