package provider

import (
	"github.com/shawntoffel/ddns/pkg/ddns"
)

type ProviderFactory struct {
	providers map[string]ddns.Provider
}

func NewProviderFactory() ddns.ProviderFactory {
	return &ProviderFactory{providers: map[string]ddns.Provider{}}
}

func (f *ProviderFactory) RegisterProvider(provider ddns.Provider) {
	if provider == nil {
		return
	}

	f.providers[provider.Name()] = provider
}

func (f *ProviderFactory) Provider(name string) ddns.Provider {
	provider, ok := f.providers[name]
	if !ok {
		return nil
	}

	return provider
}
