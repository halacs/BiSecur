package sdk

import (
	"fmt"
)

type Ports []Port

type Port struct {
	ID   int `json:"id"`
	Type int `json:"type"`
}

func (p *Port) String() string {
	portTypeStr, err := PortTypeToString(p.Type)
	if err != nil {
		portTypeStr = fmt.Sprintf("%v", err)
	}
	return fmt.Sprintf("ID=%d Type=%s", p.ID, portTypeStr)
}

func PortTypeToString(t int) (string, error) {
	switch t {
	case PORT_TYPE_NONE:
		return "NONE", nil
	case PORT_TYPE_IMPULS:
		return "IMPULS", nil
	case PORT_TYPE_AUTO_CLOSE:
		return "AUTO_CLOSE", nil
	case PORT_TYPE_ON_OFF:
		return "ON_OFF", nil
	case PORT_TYPE_UP:
		return "UP", nil
	case PORT_TYPE_DOWN:
		return "DOWN", nil
	case PORT_TYPE_HALF:
		return "HALF", nil
	case PORT_TYPE_WALK:
		return "WALK", nil
	case PORT_TYPE_LIGHT:
		return "LIGHT", nil
	case PORT_TYPE_ON:
		return "ON", nil
	case PORT_TYPE_OFF:
		return "OFF", nil
	case PORT_TYPE_LOCK:
		return "LOCK", nil
	case PORT_TYPE_UNLOCK:
		return "UNLOCK", nil
	case PORT_TYPE_OPEN_DOOR:
		return "DOOR", nil
	case PORT_TYPE_LIFT:
		return "LIFT", nil
	case PORT_TYPE_SINK:
		return "SINK", nil
	}

	return "", fmt.Errorf("unknown port type value: %d", t)
}

func (ports Ports) toString() string {
	s := ""
	for _, p := range ports {
		s = s + p.String()
	}
	return s
}
