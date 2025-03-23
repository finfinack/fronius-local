package mqtt

import (
	"encoding/json"
	"fmt"

	mq "github.com/eclipse/paho.mqtt.golang"

	"github.com/finfinack/fronius-local/pkg/config"
	"github.com/finfinack/fronius-local/pkg/fronius"
	"github.com/finfinack/logger/logging"
)

const (
	mqttQoS0 = 0 // At most once delivery
	mqttQoS1 = 1 // At least once delivery
	mqttQoS2 = 2 // Exactly once delivery
)

type MQTT struct {
	cfg    *config.MQTT
	log    *logging.Logger
	client mq.Client
}

func New(cfg *config.MQTT) (*MQTT, error) {
	m := &MQTT{
		cfg: cfg,
		log: logging.NewLogger("MQTT"),
	}

	opts := mq.NewClientOptions()
	opts.AddBroker(cfg.Broker)
	opts.SetUsername(cfg.Username)
	if cfg.Password != "" {
		opts.SetPassword(cfg.Password)
	}
	opts.OnConnect = m.onConnectHandler
	opts.OnConnectionLost = m.onConnectionLostHandler

	client := mq.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, fmt.Errorf("unable to connect MQTT: %s", token.Error())
	}
	m.client = client

	return m, nil
}

func (m *MQTT) onConnectHandler(client mq.Client) {
	m.log.Info("MQTT connected")
}

func (m *MQTT) onConnectionLostHandler(client mq.Client, err error) {
	m.log.Warnf("MQTT connection lost: %s", err)
}

func (m *MQTT) Publish(batt *fronius.Battery) error {
	message, err := json.Marshal(batt)
	if err != nil {
		return fmt.Errorf("unable to marshal message: %s", err)
	}

	token := m.client.Publish(m.cfg.Topic, mqttQoS0, false, message)
	token.Wait()
	return token.Error()
}
