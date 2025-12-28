package homeAssistant

import (
	"bisecur/cli/utils"
	"fmt"
)

func (ha *HomeAssistanceMqttClient) getDiscoveryMessage() (string, error) {
	name := ha.getUniqueObjectId()
	uniqueId := ha.getUniqueObjectId()
	device_class := `"device_class": "garage",`
	if !ha.doorStatusSupported {
		device_class = ""
	}
	commandTopic := ha.getSetPositionTopic()
	positionTopic := ha.getPositionTopicName()

	messageTemplate := `
			{
			"name": "%s",
			"unique_id": "%s",
			%s
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

	message := fmt.Sprintf(messageTemplate, name, uniqueId, device_class, commandTopic, positionTopic, uniqueId, name, ha.getAvailabilityTopic(), ha.getAvabilityMessage(true), ha.getAvabilityMessage(false))
	return message, nil
}

func (ha *HomeAssistanceMqttClient) getAvabilityMessage(online bool) string {
	if online {
		return utils.ONLINE
	} else {
		return utils.OFFLINE
	}
}
