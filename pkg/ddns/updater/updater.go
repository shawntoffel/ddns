package updater

import (
	"fmt"
	"sync"

	"github.com/rs/zerolog"
	"github.com/shawntoffel/ddns/pkg/ddns"
)

var _ ddns.Updater = &Updater{}

//Updater updates providers
type Updater struct {
	logger    zerolog.Logger
	wg        *sync.WaitGroup
	domains   []ddns.Domain
	providers map[string]ddns.Provider
}

//New returns a new Updater
func New(logger zerolog.Logger) Updater {
	return Updater{logger: logger.With().Str("component", "updater").Logger(), wg: &sync.WaitGroup{}, providers: map[string]ddns.Provider{}}
}

// RegisterProvider registers the provider
func (u *Updater) RegisterProvider(provider ddns.Provider) {
	if u.providers == nil {
		u.providers = map[string]ddns.Provider{}
	}

	u.providers[provider.Name()] = provider
}

// RegisterDomains registers domains to update
func (u *Updater) RegisterDomains(domains []ddns.Domain) {
	if u.domains == nil {
		u.domains = []ddns.Domain{}
	}

	for _, domain := range domains {
		u.domains = append(u.domains, domain)
	}
}

//Update updates all registered providers with the specified ip
func (u Updater) Update(correlationID string, ip string) {
	logger := u.logger.With().
		Str("correlationID", correlationID).
		Logger()

	for _, domain := range u.domains {
		u.wg.Add(1)

		go func(domain ddns.Domain) {
			defer u.wg.Done()

			logger.Info().
				EmbedObject(domain).
				Str("newIP", ip).
				Msg("starting update for domain")

			err := u.updateDomain(logger, domain, ip)
			if err != nil {
				logger.Error().Err(err).
					EmbedObject(domain).
					Str("newIP", ip).
					Msg("failed to update domain")

				return
			}

			logger.Info().
				EmbedObject(domain).
				Str("newIP", ip).
				Msg("finished update for domain")
		}(domain)
	}
}

// Stop gracefully stops. Waits for all spawned goroutines to finish
func (u Updater) Stop() {
	u.logger.Info().Msg("stopping")

	u.wg.Wait()
}

func (u Updater) updateDomain(logger zerolog.Logger, domain ddns.Domain, ip string) error {
	provider, ok := u.providers[domain.Provider]
	if !ok {
		return fmt.Errorf("could not find provider for domain")
	}

	records, err := provider.Records(domain)
	if err != nil {
		return fmt.Errorf("failed to get records from provider: %s", err)
	}

	for _, record := range records {
		if record.Content == ip {
			logger.Info().
				EmbedObject(domain).
				EmbedObject(record).
				Msg("record has desired content")

			continue
		}

		previousContent := record.Content
		record.Content = ip

		err = provider.UpdateRecord(record)
		if err != nil {
			logger.Error().
				Err(err).
				Str("newIP", ip).
				EmbedObject(domain).
				EmbedObject(record).
				Msg("failed to update record")

			return err
		}

		logger.Info().
			EmbedObject(record).
			EmbedObject(domain).
			Str("previousContent", previousContent).
			Msg("record has been updated")
	}

	return nil
}
