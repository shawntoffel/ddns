package provider

import (
	"strings"

	"golang.org/x/net/publicsuffix"
)

// Domain contains the domain name and records
type Domain struct {
	Name    string
	Records []string
}

// HasRecord returns true if the provided records is in the domain's slice of records
func (d Domain) HasRecord(input string) bool {
	if input == d.Name {
		input = "@"
	}

	for _, record := range d.Records {
		if record == input {
			return true
		}
	}

	return false
}

// ParseDomains parses the provided records into domains with domain records
func ParseDomains(records []string) ([]Domain, error) {
	domainMap, err := createDomainMap(records)
	if err != nil {
		return []Domain{}, err
	}

	domains := []Domain{}

	for name, records := range domainMap {
		d := Domain{
			Name:    name,
			Records: records,
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
		return ""
	}

	return record[0:pos]
}
