package homeAssistant

import (
	"bisecur/cli/homeAssistant/mockDoor"
	"bisecur/cli/utils"
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
)

const (
	qosAtLeastOnce = byte(1)
	MinTelePeriod  = 3 * time.Second
)

type HomeAssistanceMqttClient struct {
	localMac               [6]byte
	deviceMac              [6]byte
	host                   string
	port                   int
	deviceUsername         string
	devicePassword         string
	token                  uint32
	tokenCreated           time.Time
	mqttServerName         string
	mqttClientId           string
	mqttServerPort         int
	mqttServerTls          bool
	mqttServerTlsValidaton bool
	mqttBaseTopic          string
	mqttDeviceName         string
	mqttUserName           string
	mqttPassword           string
	mqttTelePeriod         time.Duration
	mqttTelePeriodFast     time.Duration
	devicePorts            []int
	log                    *logrus.Logger
	mqttClient             mqtt.Client
	requestFastUpdate      time.Time
	doorStatusSupported    bool
}

func NewHomeAssistanceMqttClient(log *logrus.Logger, localMac [6]byte, deviceMac [6]byte, deviceUsername string, devicePassword string, host string, port int, token uint32, mqttServerName string, mqttClientId string,
	mqttServerPort int, mqttServerTls bool, mqttServerTlsValidaton bool, mqttBaseTopic string,
	mqttDeviceName string, mqttUserName string, mqttPassword string, mqttTelePeriod time.Duration, mqttTelePeriodFast time.Duration, devicePorts []int, doorStatusSupported bool) (*HomeAssistanceMqttClient, error) {

	ha := &HomeAssistanceMqttClient{
		localMac:               localMac,
		deviceMac:              deviceMac,
		deviceUsername:         deviceUsername,
		devicePassword:         devicePassword,
		host:                   host,
		port:                   port,
		token:                  token,
		mqttServerName:         mqttServerName,
		mqttClientId:           mqttClientId,
		mqttServerPort:         mqttServerPort,
		mqttServerTls:          mqttServerTls,
		mqttServerTlsValidaton: mqttServerTlsValidaton,
		mqttBaseTopic:          mqttBaseTopic,
		mqttDeviceName:         mqttDeviceName,
		mqttUserName:           mqttUserName,
		mqttPassword:           mqttPassword,
		mqttTelePeriod:         mqttTelePeriod,
		mqttTelePeriodFast:     mqttTelePeriodFast,
		devicePorts:            devicePorts,
		log:                    log,
		requestFastUpdate:      time.UnixMicro(0), // initial value must be in the past
		doorStatusSupported:    doorStatusSupported,
	}

	return ha, nil
}

func (ha *HomeAssistanceMqttClient) Start() error {
	var (
		messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
			ha.log.Debugf("Received message: %s from topic: %s", msg.Payload(), msg.Topic())
		}

		homeAssistantStatusMessagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
			ha.log.Debugf("Received HA status message: %s from topic: %s", msg.Payload(), msg.Topic())
			// must not block
		}

		connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
			ha.log.Infof("Connected")
		}

		connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
			ha.log.Errorf("Connect lost: %v", err)
		}

		homeAssistantSetPossitionMessagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
			ha.log.Debugf("Received set position message: %s from topic: %s", msg.Payload(), msg.Topic())

			devicePort, err := strconv.Atoi(strings.Split(msg.Topic(), "/")[2]) // extract device port from topic
			if err != nil {
				ha.log.Errorf("failed to parse device port. %v", err)
				return
			}
			command := string(msg.Payload())

			switch command {
			case "CLOSE":
				err := ha.closeDoor(byte(devicePort))
				if err != nil {
					ha.log.Errorf("failed to close door. %v", err)
				}
			case "OPEN":
				err := ha.openDoor(byte(devicePort))
				if err != nil {
					ha.log.Errorf("failed to open door. %v", err)
				}
			case "STOP":
				err := ha.stopDoor(byte(devicePort))
				if err != nil {
					ha.log.Errorf("failed to stop door. %v", err)
				}
			case "PRESS":
				err := ha.impuls(byte(devicePort))
				if err != nil {
					ha.log.Errorf("failed to do impuls. %v", err)
				}
			}

		}
	)

	opts := mqtt.NewClientOptions()

	protocol := "tcp"
	if ha.mqttServerTls {
		protocol = "tls"
	}
	brokerUrl := fmt.Sprintf("%s://%s:%d", protocol, ha.mqttServerName, ha.mqttServerPort)
	ha.log.Debugf("MQTT Broker url: %s", brokerUrl)
	opts.AddBroker(brokerUrl)

	opts.SetClientID(ha.mqttClientId)
	opts.SetUsername(ha.mqttUserName)
	opts.SetPassword(ha.mqttPassword)
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	tlsConfig := ha.newTlsConfig()
	opts.SetTLSConfig(tlsConfig)
	opts.SetAutoReconnect(true)

	// Configure offline availability messages
	for _, devicePort := range ha.devicePorts {
		opts.SetWill(ha.getAvailabilityTopic(devicePort), ha.getAvabilityMessage(false), qosAtLeastOnce, true)
	}

	ha.mqttClient = mqtt.NewClient(opts)
	mqttToken := ha.mqttClient.Connect()
	if mqttToken.Wait() && mqttToken.Error() != nil {
		log.Fatalf("Failed to connect MQTT server. %v", mqttToken.Error())
	}
	defer func() {
		ha.log.Debugf("Disconnecting from MQTT server")

		for _, devicePort := range ha.devicePorts {
			err := ha.PublishAvabilityMessage(false, devicePort)
			if err != nil {
				ha.log.Errorf("failed to publish availability message (offline). %v", err)
			}
		}

		ha.mqttClient.Disconnect(250)

		ha.log.Infof("Disconnected from MQTT server")

		err := ha.LogoutBisecur()
		if err != nil {
			ha.log.Errorf("Error logging out of bisecur")
		} else {
			ha.log.Infof("Logged out of bisecur")
		}
	}()

	// Subscribe to home assistant's status topic (get notification when HA restarts)
	ha.mqttClient.Subscribe(utils.HomeAssistantStatusTopic, 0, homeAssistantStatusMessagePubHandler)

	for _, devicePort := range ha.devicePorts {
		// Subscribe to topics for receiving commands
		setPositionTopicName := ha.getSetPositionTopic(devicePort)
		ha.log.Debugf("Subscribing to topic: %s", setPositionTopicName)
		ha.mqttClient.Subscribe(setPositionTopicName, 0, homeAssistantSetPossitionMessagePubHandler)

		// Publish discovery messages
		err := ha.PublishDiscoveryMessage(devicePort)
		if err != nil {
			ha.log.Errorf("failed to publish discovery message. %v", err)
		}

		// Configure availabilities
		err = ha.PublishAvabilityMessage(true, int(devicePort))
		if err != nil {
			ha.log.Errorf("failed to publish availability message (online). %v", err)
		}
	}

	mockDoor.StartTicker()

	if ha.mqttTelePeriod < MinTelePeriod { // ensure minimum tele period to avoid flooding the Hormann gateway
		ha.mqttTelePeriod = MinTelePeriod
		ha.log.Warnf("Tele period is too small. Set to %v", ha.mqttTelePeriod)
	}

	if ha.mqttTelePeriodFast < MinTelePeriod { // ensure minimum tele period to avoid flooding the Hormann gateway
		ha.mqttTelePeriodFast = MinTelePeriod
		ha.log.Warnf("Tele period when door might be moving is too small. Set to %v", ha.mqttTelePeriodFast)
	}

	ticker := time.NewTicker(ha.mqttTelePeriod)
	tickerFast := time.NewTicker(ha.mqttTelePeriodFast)
	done, _ := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGINT)

out:
	for {
		select {
		case <-done.Done():
			ha.log.Infof("Exiting")
			break out
		case <-ticker.C:
			if ha.doorStatusSupported {
				for _, devicePort := range ha.devicePorts {
					err := ha.doorStatus(byte(devicePort))
					if err != nil {
						ha.log.Errorf("failed to publish current door status. %v", err)
						continue
					}
				}
			} else {
				ha.log.Debug("Status request skipped due to disable periodic requests.")
			}
		case <-tickerFast.C:
			if !ha.requestFastDootStatus() {
				continue
			}

			if ha.doorStatusSupported {
				for _, devicePort := range ha.devicePorts {
					err := ha.doorStatus(byte(devicePort))
					if err != nil {
						ha.log.Errorf("failed to publish current door status. %v", err)
						continue
					}
				}
			} else {
				ha.log.Debug("Status request skipped due to disable periodic requests.")
			}
		}
	}

	return nil
}

func (ha *HomeAssistanceMqttClient) setRequestFastDootStatus() {
	ha.requestFastUpdate = time.Now().Add(15 * time.Second) // request fast update for 15 seconds
}

func (ha *HomeAssistanceMqttClient) requestFastDootStatus() bool {
	return ha.requestFastUpdate.After(time.Now()) // check if fast update timeout is still not reached
}

func (ha *HomeAssistanceMqttClient) doorStatus(devicePort byte) error {
	ha.log.Tracef("Publish current door status")

	startTs := time.Now()
	direction, position, err := ha.getDoorStatus(devicePort)
	endTs := time.Now()
	ha.log.Debugf("Get door status took %v", endTs.Sub(startTs))
	if err != nil {
		return fmt.Errorf("failed to get door status. %v", err)
	}

	// Check if door is in motion
	if direction == utils.OPENING || direction == utils.CLOSING {
		// Maintain an increased status report frequency while the door is in motion, regardless of the trigger source.
		ha.setRequestFastDootStatus()
	}

	state := utils.UNKNOWN
	if position == 0 {
		state = utils.CLOSED
	} else if position > 0 {
		state = utils.OPEN
	}

	err = ha.PublishCurrentDoorStatus(position, direction, state, int(devicePort))
	if err != nil {
		return fmt.Errorf("failed to publish current door status. %v", err)
	}

	return nil
}

func (ha *HomeAssistanceMqttClient) PublishCurrentDoorStatus(position int, direction string, state string, devicePort int) error {
	mqttToken := ha.mqttClient.Publish(ha.getGetStateTopicName(devicePort), qosAtLeastOnce, false, state)
	if mqttToken.Wait() && mqttToken.Error() != nil {
		return fmt.Errorf("failed to publish discovery message. %v", mqttToken.Error())
	}

	mqttToken = ha.mqttClient.Publish(ha.getPositionTopicName(devicePort), qosAtLeastOnce, false, fmt.Sprintf("%d", position))
	if mqttToken.Wait() && mqttToken.Error() != nil {
		return fmt.Errorf("failed to publish discovery message. %v", mqttToken.Error())
	}

	mqttToken = ha.mqttClient.Publish(ha.getDirectionTopicName(devicePort), qosAtLeastOnce, false, direction)
	if mqttToken.Wait() && mqttToken.Error() != nil {
		return fmt.Errorf("failed to publish discovery message. %v", mqttToken.Error())
	}

	ha.log.Debugf("Published current door status: position=%d, direction=%s, state=%s", position, direction, state)

	return nil
}

func (ha *HomeAssistanceMqttClient) PublishAvabilityMessage(online bool, devicePort int) error {
	var message string

	if online {
		message = utils.ONLINE
	} else {
		message = utils.OFFLINE
	}

	mqttToken := ha.mqttClient.Publish(ha.getAvailabilityTopic(devicePort), qosAtLeastOnce, true, message)
	if mqttToken.Wait() && mqttToken.Error() != nil {
		return fmt.Errorf("failed to publish avability message. %v", mqttToken.Error())
	}

	ha.log.Debugf("Published availability message: %s", message)

	return nil
}

func (ha *HomeAssistanceMqttClient) PublishDiscoveryMessage(devicePort int) error {
	discoveryMsg, err := ha.getDiscoveryMessage(devicePort)
	if err != nil {
		return fmt.Errorf("failed to generate discovery message. %v", err)
	}

	mqttToken := ha.mqttClient.Publish(ha.getDiscoveryTopic(devicePort), qosAtLeastOnce, true, discoveryMsg)
	if mqttToken.Wait() && mqttToken.Error() != nil {
		return fmt.Errorf("failed to publish discovery message. %v", mqttToken.Error())
	}

	ha.log.Debugf("Published discovery message: %s", discoveryMsg)

	return nil
}

func (ha *HomeAssistanceMqttClient) newTlsConfig() *tls.Config {
	// https://dev.to/emqx/how-to-use-mqtt-in-golang-2oek

	certpool := x509.NewCertPool()
	/*
		ca, err := ioutil.ReadFile("ca.pem")
		if err != nil {
			log.Fatalln(err.Error())
		}
		certpool.AppendCertsFromPEM(ca)
		// Import client certificate/key pair
		clientKeyPair, err := tls.LoadX509KeyPair("client-crt.pem", "client-key.pem")
		if err != nil {
			panic(err)
		}
	*/

	return &tls.Config{
		RootCAs:            certpool,
		ClientAuth:         tls.NoClientCert,
		ClientCAs:          nil,
		InsecureSkipVerify: ha.mqttServerTlsValidaton, // #nosec #G402
		//Certificates:       []tls.Certificate{clientKeyPair},
	}
}
