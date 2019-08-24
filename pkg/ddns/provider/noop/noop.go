package noop

import (
	"github.com/shawntoffel/ddns/pkg/ddns"
)

var _ ddns.Provider = &Provider{}

// Provider is a no-op provider
type Provider struct{}

// Name returns the provider name
func (p Provider) Name() string {
	return "noop"
}

// Records returns provider records
func (p Provider) Records(domain ddns.Domain) ([]ddns.ProviderRecord, error) {
	records := []ddns.ProviderRecord{}

	for _, domainRecord := range domain.Records {
		record := ddns.ProviderRecord{
			Name:   domainRecord,
			Domain: domain.Name,
		}

		records = append(records, record)
	}

	return records, nil
}

// UpdateRecord is a noop implementation.
func (p Provider) UpdateRecord(record ddns.ProviderRecord) error {
	return nil
}
