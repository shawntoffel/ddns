package ddns

import (
	"github.com/robfig/cron"
	"github.com/rs/xid"
	"github.com/rs/zerolog"
)

// Runner checks and updates DNS records at an interval
type Runner struct {
	logger  zerolog.Logger
	checker *Checker
	updater *Updater
	knownIP string
	cron    *cron.Cron
}

//NewRunner returns a new Ddns that checks and updates DNS records at an interval
func NewRunner(logger zerolog.Logger, updater *Updater, checker *Checker) Runner {
	return Runner{logger: logger, checker: checker, updater: updater}
}

//Start starts the runner cron job with the specified interval
func (r *Runner) Start(interval string) {
	if r.cron != nil {
		r.cron.Stop()
	}

	r.cron = cron.New()
	r.cron.AddFunc("@every "+interval, r.run)
	r.cron.Start()

	r.logger.Info().Str("interval", interval).Msg("started")
}

// Stop gracefully stops
func (r *Runner) Stop() {
	r.logger.Info().Msg("stopping")

	if r.cron != nil {
		r.cron.Stop()
	}

	r.updater.Stop()
}

func (r *Runner) run() {
	correlationID := xid.New().String()

	logger := r.logger.With().
		Str("correlationID", correlationID).
		Logger()

	logger.Debug().
		Msg("checking if external IP has changed")

	externalIP, hasChanged, err := r.checker.IPHasChanged(r.knownIP)
	if err != nil {
		logger.Error().Err(err).Msg("failed to determine if external IP has changed")
		return
	}

	if !hasChanged {
		logger.Debug().
			Str("knownIP", r.knownIP).
			Str("externalIP", externalIP).
			Msg("external IP has not changed")

		return
	}

	logger.Info().
		Str("knownIP", r.knownIP).
		Str("newIP", externalIP).
		Msg("found new external IP")

	r.updater.Update(correlationID, externalIP)

	r.knownIP = externalIP
}
