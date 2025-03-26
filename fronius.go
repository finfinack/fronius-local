package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/finfinack/fronius-local/pkg/config"
	"github.com/finfinack/fronius-local/pkg/fronius"
	"github.com/finfinack/fronius-local/pkg/mqtt"
	"github.com/finfinack/logger/logging"
)

var (
	Commit  string
	Version = "dev"
)

var (
	logLevel   = flag.String("log", "INFO", "Log level")
	configFile = flag.String("config", "config.json", "File path where the configuration is stored.")
)

func fetchAndPush(f *fronius.Fronius, m *mqtt.MQTT) error {
	batt, err := f.FetchBattery()
	if err != nil {
		return fmt.Errorf("unable to fetch: %s", err)
	}
	timeOfUse, err := f.FetchTimeOfUse()
	if err != nil {
		return fmt.Errorf("unable to fetch: %s", err)
	}

	msg, err := mqtt.NewMessage(batt, timeOfUse)
	if err != nil {
		return fmt.Errorf("unable to create message: %s", err)
	}

	if err := m.Publish(msg); err != nil {
		return fmt.Errorf("unable to publish: %s", err)
	}
	return nil
}

func main() {
	flag.Parse()

	log := logging.NewLogger("MAIN")
	l, err := logging.LevelToValue(*logLevel)
	if err != nil {
		log.Fatalf("log level %q is invalid: %s\n", *logLevel, err)
	}
	logging.SetMinLogLevel(l)
	log.Infof("starting fronius-local version %q, commit %q", Version, Commit)

	for {
		log.Infof("reading config from %s", *configFile)
		cfg, err := config.ReadFromFile(*configFile)
		if err != nil {
			log.Fatalf("unable to read config: %s\n", err)
		}

		m, err := mqtt.New(cfg.MQTT)
		if err != nil {
			log.Fatalf("unable to create mqtt client: %s\n", err)
		}

		f := fronius.New(cfg.Fronius)
		if err := fetchAndPush(f, m); err != nil {
			log.Fatalf("unable to fetch and push: %s\n", err)
		}

		if cfg.Interval == 0 {
			log.Infoln("exiting")
			return
		}
		log.Infof("sleeping for %s, until %s\n", cfg.Interval, time.Now().Add(cfg.Interval).Format(time.DateTime))
		time.Sleep(time.Duration(cfg.Interval))
	}
}
