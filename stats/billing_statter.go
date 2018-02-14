package stats

import "github.com/GaruGaru/scaleway-billing-monitor/api"

type BillingStatter interface {
	Send(list api.BillingList) error
}

type StatterHolder struct {
	Statters []BillingStatter
}

func StatterHolderOf(statters ...BillingStatter) StatterHolder {
	return StatterHolder{statters}
}

func (holder StatterHolder) Send(list api.BillingList) error {
	for _, m := range holder.Statters {
		m.Send(list)
	}
	return nil
}
