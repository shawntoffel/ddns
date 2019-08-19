package ddns

import (
	"encoding/json"
	"net/http"
)

// Checker checks external IPs
type Checker struct {
	client   *http.Client
	endpoint string
}

// NewChecker returns a new Checker
func NewChecker(client *http.Client) *Checker {
	return &Checker{client: client}
}

// SetEndpoint sets the endpoint used for looking up the external IP
func (c *Checker) SetEndpoint(endpoint string) {
	c.endpoint = endpoint
}

//IPHasChanged determines if the external IP has changed. Returns the new IP if different.
func (c *Checker) IPHasChanged(knownIP string) (string, bool, error) {
	externalIP, err := c.lookupExternalIPAddress()
	if err != nil {
		return "", false, err
	}

	if externalIP == knownIP {
		return knownIP, false, nil
	}

	return externalIP, true, nil
}

type pong struct {
	Pong string
}

func (c *Checker) lookupExternalIPAddress() (string, error) {
	resp, err := c.client.Get(c.endpoint)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	pong := pong{}

	err = json.NewDecoder(resp.Body).Decode(&pong)
	if err != nil {
		return "", err
	}

	return pong.Pong, nil
}
