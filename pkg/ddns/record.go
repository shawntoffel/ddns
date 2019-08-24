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
