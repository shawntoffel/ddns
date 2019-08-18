package provider

//Provider is a dns provider
type Provider interface {
	Name() string
	SetRecords([]string)
	Update(ip string) error
}
