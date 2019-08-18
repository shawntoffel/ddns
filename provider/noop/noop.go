package noop

// Provider is a no-op provider
type Provider struct{}

// Name returns the provider name
func (p Provider) Name() string {
	return "noop"
}

// SetRecords Set the records this provider is responsible for
func (p Provider) SetRecords(records []string) error {
	return nil
}

// Update updates all records with the provided ip
func (p Provider) Update(ip string) error {
	return nil
}
