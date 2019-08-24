package ddns

import "github.com/rs/zerolog"

type ProviderRecord struct {
	ID       string
	Domain   string
	Name     string
	Type     string
	Content  string
	Metadata map[string]string
}

func (r ProviderRecord) MarshalZerologObject(e *zerolog.Event) {
	e.Str("record.ID", r.ID).
		Str("record.Domain", r.Domain).
		Str("record.Name", r.Name).
		Str("record.Type", r.Type).
		Str("record.Content", r.Content)
}

//Provider is a dns provider
type Provider interface {
	Name() string
	Records(d Domain) ([]ProviderRecord, error)
	UpdateRecord(record ProviderRecord) error
}

type ProviderFactory interface {
	RegisterProvider(provider Provider)
	Provider(name string) Provider
}

type Checker interface {
	SetEndpoint(endpoint string)
	IPHasChanged(knownIP string) (string, bool, error)
}

type Updater interface {
	AddDomains([]Domain)
	RegisterProvider(Provider)
	Update(correlationID string, ip string)
	Stop()
}

type Runner interface {
	Start(interval string)
	Stop()
}
