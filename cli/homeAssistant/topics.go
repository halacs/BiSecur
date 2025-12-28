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

func (ha *HomeAssistanceMqttClient) getDirectionTopicName() string {
	return fmt.Sprintf("%s/%s/direction", ha.mqttBaseTopic, ha.mqttDeviceName)
}

func (ha *HomeAssistanceMqttClient) getGetStateTopicName() string {
	return fmt.Sprintf("%s/%s/state", ha.mqttBaseTopic, ha.mqttDeviceName)
}

func (ha *HomeAssistanceMqttClient) getDiscoveryTopic(impulsCommandOnly bool) string {
	//<discovery_prefix>/<component>/[<node_id>/]<object_id>/config
	entity_type := "cover"
	if impulsCommandOnly {
		entity_type = "button"
	}
	return fmt.Sprintf("homeassistant/%s/halsecur/%s/config", entity_type, ha.getUniqueObjectId())
}

func (ha *HomeAssistanceMqttClient) getUniqueObjectId() string {
	deviceMacStr := hex.EncodeToString(ha.deviceMac[:])
	return fmt.Sprintf("%s%d", deviceMacStr, ha.devicePort)
}

func (ha *HomeAssistanceMqttClient) getAvailabilityTopic() string {
	return fmt.Sprintf("%s/%s/availability", ha.mqttBaseTopic, ha.mqttDeviceName)
}
