package cloudflare

import (
	"errors"

	"github.com/cloudflare/cloudflare-go"
	"github.com/shawntoffel/ddns/pkg/ddns"
)

var zoneIDMetadataKey = "zoneID"

// Provider is a Cloudflare provider
type Provider struct {
	client *cloudflare.API
}

// NewCloudflareProvider returns a new Cloudflare Provider
func NewCloudflareProvider(apiToken string) (Provider, error) {
	client, err := cloudflare.NewWithAPIToken(apiToken)
	if err != nil {
		return Provider{}, err
	}

	return Provider{client: client}, nil
}

// Name returns the provider name
func (p Provider) Name() string {
	return "cloudflare"
}

// Records returns provider records
func (p Provider) Records(domain ddns.Domain) ([]ddns.ProviderRecord, error) {
	records := []ddns.ProviderRecord{}

	zoneID, err := p.client.ZoneIDByName(domain.Name)
	if err != nil {
		return records, err
	}

	cfRecords, err := p.client.DNSRecords(zoneID, cloudflare.DNSRecord{Type: "A"})
	if err != nil {
		return records, err
	}

	for _, cfRecord := range cfRecords {
		if !domain.HasRecord(cfRecord.Name) {
			continue
		}

		record := ddns.ProviderRecord{
			ID:       cfRecord.ID,
			Name:     cfRecord.Name,
			Domain:   domain.Name,
			Type:     cfRecord.Type,
			Content:  cfRecord.Content,
			Metadata: map[string]string{zoneIDMetadataKey: zoneID},
		}

		records = append(records, record)
	}

	return records, nil
}

// UpdateRecord updates the record
func (p Provider) UpdateRecord(record ddns.ProviderRecord) error {
	zoneID, ok := record.Metadata[zoneIDMetadataKey]
	if !ok {
		return errors.New("record metadata did not contain a zoneID")
	}

	cfRecord, err := p.client.DNSRecord(zoneID, record.ID)
	if err != nil {
		return err
	}

	cfRecord.Content = record.Content

	err = p.client.UpdateDNSRecord(zoneID, record.ID, cfRecord)
	if err != nil {
		return err
	}

	return nil
}
