package ddns

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
