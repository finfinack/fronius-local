package main

import (
	"flag"

	"github.com/finfinack/fronius-local/pkg/config"
	"github.com/finfinack/fronius-local/pkg/fronius"
	"github.com/finfinack/fronius-local/pkg/mqtt"
	"github.com/finfinack/logger/logging"
)

var (
	logLevel   = flag.String("log", "INFO", "Log level")
	configFile = flag.String("config", "config.json", "File path where the configuration is stored.")
)

func main() {
	flag.Parse()

	log := logging.NewLogger("MAIN")
	l, err := logging.LevelToValue(*logLevel)
	if err != nil {
		log.Fatalf("log level %q is invalid: %s\n", *logLevel, err)
	}
	logging.SetMinLogLevel(l)

	cfg, err := config.ReadFromFile(*configFile)
	if err != nil {
		log.Fatalf("unable to read config: %s\n", err)
	}

	m, err := mqtt.New(cfg.MQTT)
	if err != nil {
		log.Fatalf("unable to create mqtt client: %s\n", err)
	}

	f := fronius.New(cfg.Fronius)
	batt, err := f.FetchBattery()
	if err != nil {
		log.Fatalf("unable to fetch: %s\n", err)
	}

	if err := m.Publish(batt); err != nil {
		log.Fatalf("unable to publish: %s\n", err)
	}
}
