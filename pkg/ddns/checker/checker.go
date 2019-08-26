package checker

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/shawntoffel/ddns/pkg/ddns"
)

var _ ddns.Checker = &Checker{}

// Checker checks external IPs
type Checker struct {
	client   *http.Client
	endpoint string
}

// New returns a new Checker with a default client & timeout
func New() Checker {
	return Checker{client: &http.Client{Timeout: 5 * time.Second}}
}

// NewWithClient returns a new Checker with the provided http client
func NewWithClient(client *http.Client) Checker {
	return Checker{client: client}
}

// SetEndpoint sets the endpoint used for looking up the external IP
func (c *Checker) SetEndpoint(endpoint string) {
	c.endpoint = endpoint
}

//IPHasChanged determines if the external IP has changed. Returns the new IP if different.
func (c Checker) IPHasChanged(knownIP string) (string, bool, error) {
	externalIP, err := c.lookupExternalIPAddress()
	if err != nil {
		return "", false, err
	}

	if externalIP == knownIP {
		return knownIP, false, nil
	}

	return externalIP, true, nil
}

func (c Checker) lookupExternalIPAddress() (string, error) {
	content, err := c.getContentFromEndpoint(c.endpoint)
	if err != nil {
		return "", err
	}

	// Many services have a trailing newline
	content = strings.TrimSuffix(content, "\n")

	externalIP := net.ParseIP(content)
	if externalIP == nil {
		return "", fmt.Errorf("endpoint did not return a valid IP")
	}

	return externalIP.String(), nil
}

func (c Checker) getContentFromEndpoint(endpoint string) (string, error) {
	resp, err := c.client.Get(endpoint)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("endpoint returned status code: %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
