package cloudflare

import (
	"github.com/cloudflare/cloudflare-go"
	"github.com/shawntoffel/ddns/provider"
)

// Provider is a Cloudflare provider
type Provider struct {
	client  *cloudflare.API
	domains []provider.Domain
}

// NewCloudflareProvider returns a new Cloudflare Provider
func NewCloudflareProvider(apiToken string) (*Provider, error) {
	client, err := cloudflare.NewWithAPIToken(apiToken)
	if err != nil {
		return &Provider{}, err
	}

	return &Provider{client: client}, nil
}

// Name returns the provider name
func (p *Provider) Name() string {
	return "cloudflare"
}

// SetDomains sets the domains this provider is responsible for
func (p *Provider) SetDomains(domains []provider.Domain) {
	p.domains = domains
}

// Update updates all records with the provided ip
func (p *Provider) Update(ip string) error {
	for _, domain := range p.domains {
		zoneID, err := p.client.ZoneIDByName(domain.Name)
		if err != nil {
			return err
		}

		cfRecords, err := p.client.DNSRecords(zoneID, cloudflare.DNSRecord{Type: "A"})
		if err != nil {
			return err
		}

		for _, cfRecord := range cfRecords {
			if domain.HasRecord(cfRecord.Name) {
				updateRecord := cfRecord
				updateRecord.Content = ip

				err = p.client.UpdateDNSRecord(zoneID, cfRecord.ID, updateRecord)
				if err != nil {
					return err
				}
			}

		}
	}

	return nil
}
