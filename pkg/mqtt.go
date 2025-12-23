package channels

import (
	"fmt"
	"time"

	mqttPaho "github.com/eclipse/paho.mqtt.golang"
	message "github.com/uug-ai/hub-pipeline-notification/message"
)

type MQTT struct {
	URI      string
	Username string
	Password string
}

func (mqtt MQTT) SendNotification(message message.Message) bool {

	MQTTURI := mqtt.URI
	MQTTUsername := mqtt.Username
	MQTTPassword := mqtt.Password

	if MQTTURI == "" {
		return false
	}

	opts := mqttPaho.NewClientOptions()

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
