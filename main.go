package main

import (
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/go-kit/kit/log"
	"github.com/jasonlvhit/gocron"
	"github.com/shawntoffel/ddnsgatewayclient"
	"os"
)

type Provider struct {
	Name    string
	Domain  string
	Records []string
}

type Config struct {
	SecondsBetweenRuns uint64
	DdnsGatewayUrl     string
	Providers          []Provider
}

func ReadConfig(fileName string, config interface{}) error {
	_, err := toml.DecodeFile(fileName, config)

	return err
}

var (
	ConfigFile string
	Logger     log.Logger
)

func init() {
	flag.StringVar(&ConfigFile, "c", "", "Config File")
	flag.Parse()
}

func UpdateDns(client ddnsgatewayclient.DdnsGatewayClient, providers []Provider) {
	updateRequest := ddnsgatewayclient.UpdateDnsRequest{}
	for _, provider := range providers {
		p := ddnsgatewayclient.Provider{}
		p.Name = provider.Name
		p.Domain = provider.Domain
		p.Records = provider.Records

		updateRequest.Providers = append(updateRequest.Providers, p)
	}

	updated, err := client.UpdateDns(updateRequest)

	if err != nil {
		Logger.Log("error", err.Error())
	}

	if updated {
		Logger.Log("dns has been updated")
	}
}

func main() {
	logger := log.NewJSONLogger(os.Stdout)
	logContext := log.With(logger, "component", "ddns", "ts", log.DefaultTimestampUTC)

	Logger = logContext

	var config = Config{}

	ReadConfig(ConfigFile, &config)

	var client = ddnsgatewayclient.NewClient(config.DdnsGatewayUrl)

	gocron.Every(config.SecondsBetweenRuns).Seconds().Do(UpdateDns, client, config.Providers)

	<-gocron.Start()

	return
}
