package ddns

import (
	"sync"

	"github.com/rs/zerolog"
	"github.com/shawntoffel/ddns/provider"
)

//Updater updates providers
type Updater struct {
	logger    zerolog.Logger
	providers []provider.Provider
	wg        *sync.WaitGroup
}

//NewUpdater returns a new Updater
func NewUpdater(logger zerolog.Logger) *Updater {
	return &Updater{logger: logger.With().Str("component", "updater").Logger(), wg: &sync.WaitGroup{}}
}

//RegisterProvider registers a provider to update
func (u *Updater) RegisterProvider(p provider.Provider) {
	if u.providers == nil {
		u.providers = []provider.Provider{}
	}

	u.providers = append(u.providers, p)

	u.logger.Info().
		Str("provider", p.Name()).
		Msg("registered provider")
}

//Update updates all registered providers with the specified ip
func (u *Updater) Update(correlationID string, ip string) {
	if u.providers == nil {
		return
	}

	logger := u.logger.With().
		Str("correlationID", correlationID).
		Logger()

	for _, p := range u.providers {
		u.updateProvider(logger, p, ip)
	}
}

// Stop gracefully stops. Waits for all spawned goroutines to finish
func (u *Updater) Stop() {
	u.logger.Info().Msg("stopping")

	u.wg.Wait()
}

func (u *Updater) updateProvider(logger zerolog.Logger, p provider.Provider, ip string) {
	u.wg.Add(1)

	go func() {
		defer u.wg.Done()

		logger.Info().
			Str("provider", p.Name()).
			Str("newIP", ip).
			Msg("starting update for provider")

		err := p.Update(ip)
		if err != nil {
			logger.Error().
				Err(err).
				Str("provider", p.Name()).
				Str("newIP", ip).
				Msg("failed to update provider")
		}

		logger.Info().
			Str("provider", p.Name()).
			Str("newIP", ip).
			Msg("finished update for provider")
	}()
}
