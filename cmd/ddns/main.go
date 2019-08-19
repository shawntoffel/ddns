package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/shawntoffel/ddns"
	"github.com/shawntoffel/ddns/provider"
	"github.com/shawntoffel/ddns/provider/cloudflare"
	"github.com/shawntoffel/ddns/provider/digitalocean"
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

	updater := ddns.NewUpdater(logger)

	if flagNoopRecords != "" {
		domains, err := provider.ParseDomains(strings.Split(flagNoopRecords, ","))
		if err != nil {
			logger.Error().Err(err).Msg("failed to parse noop domains")
			os.Exit(1)
		}

		noopProvider := &noop.Provider{}
		noopProvider.SetDomains(domains)

		updater.RegisterProvider(noopProvider)
	}

	if flagDigitalOceanRecords != "" {
		apiToken, err := readSecretFromFile(os.Getenv(digitalOceanTokenFile))
		if err != nil {
			logger.Error().Err(err).Msg("failed to read digital ocean api token")
			os.Exit(1)
		}

		domains, err := provider.ParseDomains(strings.Split(flagDigitalOceanRecords, ","))
		if err != nil {
			logger.Error().Err(err).Msg("failed to parse digital ocean domains")
			os.Exit(1)
		}

		digitaloceanProvider := digitalocean.NewDigitalOceanProvider(apiToken)
		digitaloceanProvider.SetDomains(domains)

		updater.RegisterProvider(digitaloceanProvider)
	}

	if flagCloudflareRecords != "" {
		apiKey, err := readSecretFromFile(os.Getenv(cloudFlareKeyFile))
		if err != nil {
			logger.Error().Err(err).Msg("failed to read cloudflare api key")
			os.Exit(1)
		}

		email, err := readSecretFromFile(os.Getenv(cloudFlareEmailFile))
		if err != nil {
			logger.Error().Err(err).Msg("failed to read cloudflare email")
			os.Exit(1)
		}

		domains, err := provider.ParseDomains(strings.Split(flagCloudflareRecords, ","))
		if err != nil {
			logger.Error().Err(err).Msg("failed to parse digital ocean domains")
			os.Exit(1)
		}

		cloudflareProvider, err := cloudflare.NewCloudflareProvider(apiKey, email)
		if err != nil {
			logger.Error().Err(err).Msg("failed to initialize cloudflare provider")
			os.Exit(1)
		}

		cloudflareProvider.SetDomains(domains)

		updater.RegisterProvider(cloudflareProvider)
	}

	checker := ddns.NewChecker(&http.Client{})
	checker.SetEndpoint(flagEndpoint)

	runner := ddns.NewRunner(logger, updater, checker)
	runner.Start(flagInterval)

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGINT)

	select {
	case sig := <-sigChan:
		runner.Stop()
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
