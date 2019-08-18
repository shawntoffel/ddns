package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/shawntoffel/ddns"
	"github.com/shawntoffel/ddns/provider/noop"
)

// Version of ddns
var Version = ""

const (
	digitalOceanTokenFile = "DDNS_DIGITALOCEAN_TOKEN_FILE"
	cloudFlareKeyFile     = "DDNS_CLOUDFLARE_KEY_FILE"
	cloudFlareEmailFile   = "DDNS_CLOUDFLARE_EMAIL_FILE"
)

var (
	flagVersion     = false
	flagDebug       = false
	flagInterval    = "30s"
	flagEndpoint    = "http://ping.shawntoffel.com:10002/ping"
	flagNoopRecords = ""
)

func parseCli() {
	flag.BoolVar(&flagVersion, "v", false, "version")
	flag.BoolVar(&flagDebug, "D", false, "enable debug logging")
	flag.StringVar(&flagInterval, "i", flagInterval, "run interval")
	flag.StringVar(&flagEndpoint, "endpoint", flagEndpoint, "endpoint used to lookup the external IP")
	flag.StringVar(&flagNoopRecords, "records.noop", flagNoopRecords, "comma delimited list of noop records")

	flag.Parse()
}

func main() {
	parseCli()

	if flagVersion {
		fmt.Println(Version)
		os.Exit(0)
	}

	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if flagDebug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	logger := zerolog.New(os.Stdout).With().Timestamp().
		Str("version", Version).
		Logger()

	updater := ddns.NewUpdater(logger.With().Str("component", "updater").Logger())

	if flagNoopRecords != "" {
		noopProvider := &noop.Noop{}
		noopProvider.SetRecords(strings.Split(flagNoopRecords, ","))

		updater.RegisterProvider(noopProvider)
	}

	checker := ddns.NewChecker(&http.Client{})
	checker.SetEndpoint(flagEndpoint)

	runner := ddns.NewRunner(logger.With().Str("component", "runner").Logger(), updater, checker)

	runner.Start(flagInterval)

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGINT)

	select {
	case sig := <-sigChan:
		runner.Stop()
		logger.Info().Err(fmt.Errorf("%s", sig)).Msg("stopped")
	}
}
