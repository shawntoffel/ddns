package noop

// Noop is a no-op provider
type Noop struct {
	records []string
}

// Name returns the provider name
func (n *Noop) Name() string {
	return "noop"
}

// SetRecords Set the records this provider is responsible for
func (n *Noop) SetRecords(records []string) {
	n.records = records
}

// Update updates all records with the provided ip
func (n *Noop) Update(ip string) error {
	return nil
}
