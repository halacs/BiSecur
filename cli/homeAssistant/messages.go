package homeAssistant

import (
	"bisecur/cli/utils"
	"fmt"
)

func (ha *HomeAssistanceMqttClient) getDiscoveryMessage() (string, error) {
	name := ha.getUniqueObjectId()
	uniqueId := ha.getUniqueObjectId()
	commandTopic := ha.getSetPositionTopic()
	positionTopic := ha.getPositionTopicName()

	messageTemplate := `
			{
			"name": "%s",
			"unique_id": "%s",
			"device_class": "garage",
			"command_topic": "%s",
			"position_topic": "%s",
			"device": {
    			"identifiers": ["%s"],
			    "name": "%s"
  				},
			"availability_topic": "%s",
			"payload_available": "%s",
			"payload_not_available": "%s"
			}`

	message := fmt.Sprintf(messageTemplate, name, uniqueId, commandTopic, positionTopic, uniqueId, name, ha.getAvailabilityTopic(), ha.getAvabilityMessage(true), ha.getAvabilityMessage(false))
	return message, nil
}

func (ha *HomeAssistanceMqttClient) getAvabilityMessage(online bool) string {
	if online {
		return utils.ONLINE
	} else {
		return utils.OFFLINE
	}
}
