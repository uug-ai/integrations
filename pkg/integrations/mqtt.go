package integrations

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	mqttPaho "github.com/eclipse/paho.mqtt.golang"
	"github.com/uug-ai/integrations/pkg/message"
)

// MQTTOptions configures MQTT.
type MQTTOptions struct {
	// URI is the broker URI.
	URI string `json:"uri,omitempty" validate:"required"`
	// Username is the broker user.
	Username string `json:"username,omitempty"`
	// Password is the broker password.
	Password string `json:"password,omitempty"`
}

// MQTTOptionsBuilder builds MQTTOptions.
type MQTTOptionsBuilder struct {
	options *MQTTOptions
}

// NewMQTTOptions starts building MQTTOptions.
func NewMQTTOptions() *MQTTOptionsBuilder {
	return &MQTTOptionsBuilder{options: &MQTTOptions{}}
}

// SetURI sets the broker URI.
func (b *MQTTOptionsBuilder) SetURI(uri string) *MQTTOptionsBuilder {
	b.options.URI = uri
	return b
}

// SetUsername sets the broker user.
func (b *MQTTOptionsBuilder) SetUsername(username string) *MQTTOptionsBuilder {
	b.options.Username = username
	return b
}

// SetPassword sets the broker password.
func (b *MQTTOptionsBuilder) SetPassword(password string) *MQTTOptionsBuilder {
	b.options.Password = password
	return b
}

// Build returns the configured MQTTOptions.
func (b *MQTTOptionsBuilder) Build() *MQTTOptions {
	return b.options
}

// MQTT connects to an MQTT broker.
type MQTT struct {
	options *MQTTOptions
}

// NewMQTT validates opts and returns a MQTT sender that uses a copy of them.
func NewMQTT(opts *MQTTOptions) (*MQTT, error) {
	if err := validateOptions("mqtt", opts); err != nil {
		return nil, err
	}
	options := *opts
	return &MQTT{options: &options}, nil
}

func newMQTTClientID() (string, error) {
	return newMQTTClientIDWithPrefix("hub-pipe-notif-")
}

func newMQTTClientIDWithPrefix(prefix string) (string, error) {
	random := make([]byte, 4)
	if _, err := rand.Read(random); err != nil {
		return "", err
	}

	return prefix + hex.EncodeToString(random), nil
}

// mqttPublishTimeout bounds connecting and having a message acknowledged
// in Publish.
var mqttPublishTimeout = 10 * time.Second

// Publish connects to the broker, publishes payload on topic with QoS 1
// (not retained), waits for the broker to acknowledge it, and
// disconnects. It returns an error when the broker cannot be reached or
// does not acknowledge the message in time, so the caller can retry or
// report it.
func (mqtt *MQTT) Publish(topic string, payload []byte) error {
	if strings.TrimSpace(topic) == "" {
		return errors.New("mqtt: topic is empty")
	}
	clientID, err := newMQTTClientIDWithPrefix("hub-publish-")
	if err != nil {
		return fmt.Errorf("mqtt: client id: %w", err)
	}
	opts := mqttPaho.NewClientOptions()
	opts.SetClientID(clientID)
	opts.AddBroker(mqtt.options.URI)
	if mqtt.options.Username != "" || mqtt.options.Password != "" {
		opts.SetUsername(mqtt.options.Username)
		opts.SetPassword(mqtt.options.Password)
	}
	opts.SetCleanSession(true)
	opts.SetConnectRetry(false)
	opts.SetAutoReconnect(false)
	opts.SetConnectTimeout(mqttPublishTimeout)

	client := mqttPaho.NewClient(opts)
	connect := client.Connect()
	if !connect.WaitTimeout(mqttPublishTimeout) {
		return errors.New("mqtt: timed out connecting to the broker")
	}
	if err := connect.Error(); err != nil {
		return fmt.Errorf("mqtt: connect: %w", err)
	}
	defer client.Disconnect(250)

	sent := client.Publish(topic, 1, false, payload)
	if !sent.WaitTimeout(mqttPublishTimeout) {
		return errors.New("mqtt: timed out waiting for the broker to acknowledge the message")
	}
	if err := sent.Error(); err != nil {
		return fmt.Errorf("mqtt: publish: %w", err)
	}
	return nil
}

func (mqtt *MQTT) SendNotification(message message.Message) bool {

	MQTTURI := mqtt.options.URI
	MQTTUsername := mqtt.options.Username
	MQTTPassword := mqtt.options.Password

	if MQTTURI == "" {
		return false
	}

	opts := mqttPaho.NewClientOptions()
	clientID, err := newMQTTClientID()
	if err != nil {
		fmt.Println("Unable to generate MQTT client ID")
		return false
	}
	opts.SetClientID(clientID)

	// We will set the MQTT endpoint to which we want to connect
	// and share and receive messages to/from.
	opts.AddBroker(MQTTURI)

	// Our MQTT broker can have username/password credentials
	// to protect it from the outside.
	if MQTTUsername != "" && MQTTPassword != "" {
		opts.SetUsername(MQTTUsername)
		opts.SetPassword(MQTTPassword)
	}

	// Some extra options to make sure the connection behaves
	// properly. More information here: github.com/eclipse/paho.mqtt.golang.
	opts.SetCleanSession(true)
	opts.SetConnectRetry(true)
	//opts.SetAutoReconnect(true)
	opts.SetConnectTimeout(30 * time.Second)

	opts.OnConnect = func(c mqttPaho.Client) {
		// We managed to connect to the MQTT broker, hurray!
		//log.Log.Info("routers.mqtt.main.ConfigureMQTT(): " + mqttClientID + " connect
		fmt.Println("Connected to MQTT broker")
	}

	mqc := mqttPaho.NewClient(opts)
	if token := mqc.Connect(); token.WaitTimeout(3 * time.Second) {
		if token.Error() != nil {
			//log.Log.Error("routers.mqtt.main.ConfigureMQTT(): unable to establish mqtt broker connection, error was: " + token.Error().Error())
			fmt.Println("Unable to establish mqtt broker connection")
		}
	}

	/*posString := fmt.Sprintf("%f,%f,%f", pos.PanTilt.X, pos.PanTilt.Y, pos.Zoom.X)
	m := models.Message{
		Payload: models.Payload{
			Action:   "ptz-position",
			DeviceId: configuration.Config.Key,
			Value: map[string]interface{}{
				"timestamp": positionPayload.Timestamp,
				"position":  posString,
			},
		},
	}
	payload, err := models.PackageMQTTMessage(configuration, m)
	if err == nil {
		mqttClient.Publish("kerberos/hub/"+hubKey, 0, false, payload)
	} else {
		fmt.Println("Something went wrong while sending position to hub")
		//log.Log.Info("routers.mqtt.main.HandlePTZPosition(): something went wrong while sending position to hub: " + string(payload))
	}*/

	return true
}
