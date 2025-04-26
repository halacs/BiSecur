package homeAssistant

import (
	"bisecur/cli"
	"bisecur/cli/bisecur"
	"bisecur/cli/utils"
	"fmt"
	"time"
)

func (ha *HomeAssistanceMqttClient) autoLogin() error {
	if ha.tokenCreated.Add(bisecur.TokenExpirationTime).Before(time.Now()) {
		cli.Log.Infof("Token expired. Logging in...")

		var err error
		ha.token, err = bisecur.Login(ha.localMac, ha.deviceMac, ha.host, ha.port, ha.deviceUsername, ha.devicePassword)
		if err != nil {
			return fmt.Errorf("failed to auto login. %v", err)
		}

		ha.tokenCreated = time.Now() // note when token was received
	}

	return nil
}

func (ha *HomeAssistanceMqttClient) getDiscoveryMessage() (string, error) {
	name := ha.getUniqueObjectId()
	uniqueId := name
	commandTopic := "halsecur/cmnd/garage/position"
	positionTopic := "halsecur/garage/position"
	tiltStatusTopic := "halsecur/garage/tilt_status"

	messageTemplate := `
			{
			"name": "%s",
			"unique_id": "%s",
			"device_class": "garage",
			"command_topic": "%s",
			"position_topic": "%s",
			"tilt_status_topic": "%s",
			"device": {
    			"identifiers": ["%s"],
			    "name": "%s"
  				},
			"availability_topic": "%s",
			"payload_available": "%s",
			"payload_not_available": "%s"
			}`

	message := fmt.Sprintf(messageTemplate, name, uniqueId, commandTopic, positionTopic, tiltStatusTopic, uniqueId, name, ha.getAvailabilityTopic(), ha.getAvabilityMessage(true), ha.getAvabilityMessage(false))
	return message, nil
}

func (ha *HomeAssistanceMqttClient) getAvabilityMessage(online bool) string {
	if online {
		return utils.ONLINE
	} else {
		return utils.OFFLINE
	}
}
