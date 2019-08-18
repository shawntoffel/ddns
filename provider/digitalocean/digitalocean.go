package digitalocean

import (
	"context"
	"strings"

	"github.com/digitalocean/godo"
	"golang.org/x/net/publicsuffix"
	"golang.org/x/oauth2"
)

// Provider is a Digital Ocean provider
type Provider struct {
	client  *godo.Client
	domains map[string][]string
}

// NewDigitalOceanProvider returns a new DigitalOceanProvider
func NewDigitalOceanProvider(apiToken string) Provider {
	tokenSource := &tokenSource{
		AccessToken: apiToken,
	}

	oauthClient := oauth2.NewClient(context.Background(), tokenSource)
	doClient := godo.NewClient(oauthClient)

	return Provider{client: doClient, domains: map[string][]string{}}
}

// SetRecords Set the records this provider is responsible for
func (p Provider) SetRecords(records []string) error {
	for _, record := range records {
		domain, err := publicsuffix.EffectiveTLDPlusOne(record)
		if err != nil {
			return err
		}

		subdomain := p.parseSubdomain(record, domain)

		p.domains[domain] = append(p.domains[domain], subdomain)
	}

	return nil
}

// Name returns the provider name
func (p Provider) Name() string {
	return "digitalocean"
}

// Update updates all records with the provided ip
func (p Provider) Update(ip string) error {
	for domain := range p.domains {
		err := p.updateDomain(domain, ip)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p Provider) updateDomain(domain string, ip string) error {
	doRecords, _, err := p.client.Domains.Records(context.Background(), domain, nil)
	if err != nil {
		return err
	}

	subdomains := p.domains[domain]

	for _, doRecord := range doRecords {
		if doRecord.Type != "A" {
			continue
		}

		for _, record := range subdomains {
			if doRecord.Name == record {
				edit := &godo.DomainRecordEditRequest{
					Data: ip,
				}

				_, _, err := p.client.Domains.EditRecord(context.Background(), domain, doRecord.ID, edit)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (p Provider) parseSubdomain(record string, domain string) string {
	if record == domain {
		return "@"
	}

	pos := strings.Index(record, "."+domain)
	if pos == -1 {
		return ""
	}

	return record[0:pos]
}
