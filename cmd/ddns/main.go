package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/shawntoffel/ddns/pkg/ddns"
	"github.com/shawntoffel/ddns/pkg/ddns/checker"
	"github.com/shawntoffel/ddns/pkg/ddns/provider/cloudflare"
	"github.com/shawntoffel/ddns/pkg/ddns/provider/digitalocean"
	"github.com/shawntoffel/ddns/pkg/ddns/provider/noop"
	"github.com/shawntoffel/ddns/pkg/ddns/runner"
	"github.com/shawntoffel/ddns/pkg/ddns/updater"
)

// Version of ddns
var Version = ""

const (
	digitalOceanTokenFile = "DDNS_DIGITALOCEAN_TOKEN_FILE"
	cloudFlareTokenFile   = "DDNS_CLOUDFLARE_TOKEN_FILE"
)

var (
	flagVersion             = false
	flagDebug               = false
	flagInterval            = "30s"
	flagEndpoint            = "http://ping.shawntoffel.com:10002/ping"
	flagNoopRecords         = ""
	flagDigitalOceanRecords = ""
	flagCloudflareRecords   = ""
)

func parseCli() {
	flag.BoolVar(&flagVersion, "v", false, "version")
	flag.BoolVar(&flagDebug, "D", false, "enable debug logging")
	flag.StringVar(&flagInterval, "i", flagInterval, "run interval")
	flag.StringVar(&flagEndpoint, "endpoint", flagEndpoint, "endpoint used to lookup the external IP")
	flag.StringVar(&flagNoopRecords, "records.noop", flagNoopRecords, "comma delimited list of noop records")
	flag.StringVar(&flagDigitalOceanRecords, "records.digitalocean", flagDigitalOceanRecords, "comma delimited list of digitalocean records")
	flag.StringVar(&flagCloudflareRecords, "records.cloudflare", flagCloudflareRecords, "comma delimited list of cloudflare records")

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

	u := updater.New(logger)

	if flagNoopRecords != "" {
		provider := &noop.Provider{}
		u.RegisterProvider(provider)

		domains, err := ddns.ParseDomains(strings.Split(flagNoopRecords, ","), provider.Name())
		if err != nil {
			logger.Error().Err(err).Msg("failed to parse noop domains")
			os.Exit(1)
		}

		u.RegisterDomains(domains)
	}

	if flagDigitalOceanRecords != "" {
		apiToken, err := readSecretFromFile(os.Getenv(digitalOceanTokenFile))
		if err != nil {
			logger.Error().Err(err).Msg("failed to read digitalocean api token")
			os.Exit(1)
		}

		provider := digitalocean.NewDigitalOceanProvider(apiToken)
		u.RegisterProvider(&provider)

		domains, err := ddns.ParseDomains(strings.Split(flagDigitalOceanRecords, ","), provider.Name())
		if err != nil {
			logger.Error().Err(err).Msg("failed to parse digitalocean domains")
			os.Exit(1)
		}

		u.RegisterDomains(domains)
	}

	if flagCloudflareRecords != "" {
		apiToken, err := readSecretFromFile(os.Getenv(cloudFlareTokenFile))
		if err != nil {
			logger.Error().Err(err).Msg("failed to read cloudflare api key")
			os.Exit(1)
		}

		provider, err := cloudflare.NewCloudflareProvider(apiToken)
		if err != nil {
			logger.Error().Err(err).Msg("failed to initialize cloudflare provider")
			os.Exit(1)
		}

		u.RegisterProvider(&provider)

		domains, err := ddns.ParseDomains(strings.Split(flagCloudflareRecords, ","), provider.Name())
		if err != nil {
			logger.Error().Err(err).Msg("failed to parse cloudflare domains")
			os.Exit(1)
		}

		u.RegisterDomains(domains)
	}

	c := checker.New()
	c.SetEndpoint(flagEndpoint)

	r := runner.New(logger, &u, &c)
	r.Start(flagInterval)

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGTERM)
	signal.Notify(sigChan, syscall.SIGINT)

	select {
	case sig := <-sigChan:
		r.Stop()
		logger.Info().Err(fmt.Errorf("%s", sig)).Msg("stopped")
	}
}

func readSecretFromFile(filename string) (string, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", errors.New("could not read secrets file: " + filename + ": " + err.Error())
	}

	return string(b), nil
}
