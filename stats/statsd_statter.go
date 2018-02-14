package stats

import (
	"github.com/GaruGaru/scaleway-billing-monitor/api"
	"strconv"
	"github.com/cactus/go-statsd-client/statsd"
	"errors"
	"os"
)

type StatsdStatter struct {
	statsd statsd.Statter
}

func NewStatsdStatterWithFallback() StatsdStatter {
	var client statsd.Statter
	var err error
	client, err = statsd.NewClient(os.Getenv("STATSD_HOST"), os.Getenv("STATSD_PREFIX"))
	if err != nil {
		client, err = statsd.NewNoopClient()
	}
	return StatsdStatter{statsd: client}
}

func NewStatsdStatter(host string, prefix string) (StatsdStatter, error) {
	client, err := statsd.NewClient(host, prefix)

	if err != nil {
		return StatsdStatter{}, err
	}

	return StatsdStatter{
		statsd: client,
	}, nil

}

func (statter StatsdStatter) Send(list api.BillingList) error {
	if len(list.Invoices) > 0 {
		invoice := list.Invoices[0]

		totalTaxed, _ := fstringToInt64(invoice.TotalTaxed)
		totalTax, _ := fstringToInt64(invoice.TotalTax)
		totalUntaxed, _ := fstringToInt64(invoice.TotalUntaxed)
		totalUndiscounted, _ := fstringToInt64(invoice.TotalUndiscounted)

		valuesMap := make(map[string]int64)
		valuesMap["current.total-taxed"] = totalTaxed
		valuesMap["current.total-tax"] = totalTax
		valuesMap["current.total-untaxed"] = totalUntaxed
		valuesMap["current.total-undiscounted"] = totalUndiscounted

		for k, v := range valuesMap {
			statter.statsd.Gauge(k, v, 1.0)
		}

	}

	return nil
}

func fstringToInt64(in string) (int64, error) {
	floatValue, err := strconv.ParseFloat(in, 32)
	if err != nil {
		return 0, errors.New("Invalid float string " + in)
	}
	return int64(floatValue * 100), nil
}
