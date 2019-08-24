package digitalocean

import (
	"context"
	"strconv"

	"github.com/digitalocean/godo"
	"github.com/shawntoffel/ddns/pkg/ddns"
	"golang.org/x/oauth2"
)

var _ ddns.Provider = &Provider{}

// Provider is a Digital Ocean provider
type Provider struct {
	client *godo.Client
}

// NewDigitalOceanProvider returns a new DigitalOcean Provider
func NewDigitalOceanProvider(apiToken string) Provider {
	tokenSource := &tokenSource{
		AccessToken: apiToken,
	}

	oauthClient := oauth2.NewClient(context.Background(), tokenSource)
	doClient := godo.NewClient(oauthClient)

	return Provider{client: doClient}
}

// Name returns the provider name
func (p Provider) Name() string {
	return "digitalocean"
}

// Records returns provider records
func (p Provider) Records(domain ddns.Domain) ([]ddns.ProviderRecord, error) {
	records := []ddns.ProviderRecord{}

	doRecords, _, err := p.client.Domains.Records(context.Background(), domain.Name, nil)
	if err != nil {
		return records, err
	}

	for _, doRecord := range doRecords {
		if doRecord.Type != "A" {
			continue
		}

		if !domain.HasRecord(doRecord.Name) {
			continue
		}

		record := ddns.ProviderRecord{
			ID:      strconv.Itoa(doRecord.ID),
			Name:    doRecord.Name,
			Domain:  domain.Name,
			Type:    doRecord.Type,
			Content: doRecord.Data,
		}

		records = append(records, record)
	}

	return records, nil
}

// UpdateRecord updates the record
func (p Provider) UpdateRecord(record ddns.ProviderRecord) error {
	edit := &godo.DomainRecordEditRequest{
		Data: record.Content,
	}

	id, err := strconv.Atoi(record.ID)
	if err != nil {
		return err
	}
	_, _, err = p.client.Domains.EditRecord(context.Background(), record.Domain, id, edit)
	if err != nil {
		return err
	}

	return nil
}
