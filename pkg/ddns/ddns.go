package ddns

import "github.com/rs/zerolog"

// ProviderRecord represents a record from a provider
type ProviderRecord struct {
	ID       string
	Domain   string
	Name     string
	Type     string
	Content  string
	Metadata map[string]interface{}
}

// MarshalZerologObject marshals a ProviderRecord for logging
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

// Checker checks external IPs
type Checker interface {
	SetEndpoint(endpoint string)
	IPHasChanged(knownIP string) (string, bool, error)
}

//Updater updates providers
type Updater interface {
	RegisterDomains([]Domain)
	RegisterProvider(Provider)
	Update(correlationID string, ip string)
	Stop()
}

// Runner checks and updates DNS records at an interval
type Runner interface {
	Start(interval string)
	Stop()
}
