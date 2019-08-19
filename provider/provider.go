package provider

//Provider is a dns provider
type Provider interface {
	Name() string
	SetDomains([]Domain)
	Update(ip string) error
}
