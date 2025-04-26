package homeAssistant

import (
	"encoding/hex"
	"fmt"
)

func (ha *HomeAssistanceMqttClient) getPositionTopicName() string {
	return fmt.Sprintf("%s/%s/position", ha.mqttBaseTopic, ha.mqttDeviceName)
}

func (ha *HomeAssistanceMqttClient) getSetPositionTopic() string {
	return fmt.Sprintf("%s/cmnd/%s/position", ha.mqttBaseTopic, ha.mqttDeviceName)
}

func (ha *HomeAssistanceMqttClient) getGetStateTopicName() string {
	return fmt.Sprintf("%s/%s/state", ha.mqttBaseTopic, ha.mqttDeviceName)
}

func (ha *HomeAssistanceMqttClient) getDiscoveryTopic() string {
	//<discovery_prefix>/<component>/[<node_id>/]<object_id>/config

	return fmt.Sprintf("homeassistant/cover/halsecur/%s/config", ha.getUniqueObjectId())
}

func (ha *HomeAssistanceMqttClient) getUniqueObjectId() string {
	deviceMacStr := hex.EncodeToString(ha.deviceMac[:])
	return fmt.Sprintf("%s%d", deviceMacStr, ha.devicePort)
}

func (ha *HomeAssistanceMqttClient) getAvailabilityTopic() string {
	return fmt.Sprintf("%s/%s/availability", ha.mqttBaseTopic, ha.mqttDeviceName)
}
