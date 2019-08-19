package digitalocean

import (
	"context"

	"github.com/digitalocean/godo"
	"github.com/shawntoffel/ddns/provider"
	"golang.org/x/oauth2"
)

// Provider is a Digital Ocean provider
type Provider struct {
	client  *godo.Client
	domains []provider.Domain
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

// SetDomains sets the domains this provider is responsible for
func (p Provider) SetDomains(domains []provider.Domain) {
	p.domains = domains
}

// Update updates all domain records with the provided ip
func (p Provider) Update(ip string) error {
	for _, domain := range p.domains {
		doRecords, _, err := p.client.Domains.Records(context.Background(), domain.Name, nil)
		if err != nil {
			return err
		}

		for _, doRecord := range doRecords {
			if doRecord.Type != "A" {
				continue
			}

			if domain.HasRecord(doRecord.Name) {
				edit := &godo.DomainRecordEditRequest{
					Data: ip,
				}

				_, _, err := p.client.Domains.EditRecord(context.Background(), domain.Name, doRecord.ID, edit)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
