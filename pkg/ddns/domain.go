package ddns

import (
	"strings"

	"github.com/rs/zerolog"
	"golang.org/x/net/publicsuffix"
)

// Domain holds domain info
type Domain struct {
	Name     string
	Records  []string
	Provider string
}

// MarshalZerologObject marshals a Domain for logging
func (d Domain) MarshalZerologObject(e *zerolog.Event) {
	e.Str("domain.Name", d.Name).
		Str("domain.Records", strings.Join(d.Records, ", ")).
		Str("domain.Provider", d.Provider)
}

// HasRecord returns true if the provided records is in the domain's slice of records
func (d Domain) HasRecord(input string) bool {
	input = parseDomainRecord(input, d.Name)

	for _, record := range d.Records {
		if record == input {
			return true
		}
	}

	return false
}

// ParseDomains parses a slice of string domains into Domain types
func ParseDomains(records []string, provider string) ([]Domain, error) {
	domainMap, err := createDomainMap(records)
	if err != nil {
		return []Domain{}, err
	}

	domains := []Domain{}

	for name, records := range domainMap {
		d := Domain{
			Name:     name,
			Records:  records,
			Provider: provider,
		}

		domains = append(domains, d)
	}

	return domains, nil
}

func createDomainMap(records []string) (map[string][]string, error) {
	domainMap := map[string][]string{}

	for _, record := range records {
		domain, err := parseDomain(record)
		if err != nil {
			return domainMap, err
		}

		domainRecord := parseDomainRecord(record, domain)

		domainMap[domain] = append(domainMap[domain], domainRecord)
	}

	return domainMap, nil
}

func parseDomain(record string) (string, error) {
	domain, err := publicsuffix.EffectiveTLDPlusOne(record)
	if err != nil {
		return "", err
	}

	return domain, nil
}

func parseDomainRecord(record string, domain string) string {
	if record == domain {
		return "@"
	}

	pos := strings.Index(record, "."+domain)
	if pos == -1 {
		return record
	}

	return record[0:pos]
}
