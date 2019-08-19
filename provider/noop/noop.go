package noop

import "github.com/shawntoffel/ddns/provider"

// Provider is a no-op provider
type Provider struct{}

// Name returns the provider name
func (p *Provider) Name() string {
	return "noop"
}

// SetDomains Set the domains this provider is responsible for
func (p *Provider) SetDomains(domains []provider.Domain) {}

// Update updates all records with the provided ip
func (p *Provider) Update(ip string) error {
	return nil
}
