package mqtt

import (
	"encoding/json"
	"fmt"
	"net"

	"dumper/internal/model"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/zap"
)

type MQTT struct {
	Logger *zap.Logger
	Client mqtt.Client
	Config *Config
}

func NewMQTT(logger *zap.Logger, config *Config) *MQTT {
	// Создание клиента MQTT
	onMessageReceivedCallback := onMessageReceived(logger)
	opts := mqtt.NewClientOptions()
	// // адрес и порт RabbitMQ MQTT-брокера
	opts.AddBroker(fmt.Sprintf("tcp://%s", net.JoinHostPort(config.Host, config.Port)))
	opts.SetClientID(config.ClientID)
	opts.Username = config.Username
	opts.Password = config.Password
	opts.SetDefaultPublishHandler(onMessageReceivedCallback)
	client := mqtt.NewClient(opts)

	return &MQTT{
		Logger: logger,
		Client: client,
		Config: config,
	}
}

func onMessageReceived(logger *zap.Logger) func(client mqtt.Client, message mqtt.Message) {
	return func(client mqtt.Client, message mqtt.Message) {
		logger.Info(
			fmt.Sprintf("Received message on topic: %s", message.Topic()))
		logger.Info(
			fmt.Sprintf("Message: %s", string(message.Payload())))

		m := model.Message{}
		err := json.Unmarshal(message.Payload(), &m)
		if err != nil {
			logger.Error("unmarshal error: ", zap.Error(err))
		}
	}
}

func (m *MQTT) Stop() error {
	return nil
}

func (m *MQTT) Start() error {
	// Подключение к MQTT-брокеру
	if token := m.Client.Connect(); token.Wait() && token.Error() != nil {
		return fmt.Errorf("%w", token.Error())
	}

	onMessageReceivedCallback := onMessageReceived(m.Logger)

	topic := m.Config.Topic
	if token := m.Client.Subscribe(topic, 0, onMessageReceivedCallback); token.Wait() && token.Error() != nil {
		return fmt.Errorf("%w", token.Error())
	}

	return nil
}
