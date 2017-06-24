package main

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/jasonlvhit/gocron"
	"github.com/shawntoffel/ddnsgatewayclient"
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

	_, err := client.UpdateDns(updateRequest)

	if err != nil {
		fmt.Printf(err.Error())
	}
}

func main() {

	var config = Config{}

	ReadConfig(ConfigFile, &config)

	var client = ddnsgatewayclient.NewClient(config.DdnsGatewayUrl)

	gocron.Every(config.SecondsBetweenRuns).Seconds().Do(UpdateDns, client, config.Providers)

	<-gocron.Start()

	return
}
