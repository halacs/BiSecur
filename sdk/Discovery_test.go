package sdk

import (
	"encoding/xml"
	"fmt"
	"reflect"
	"testing"
)

func TestUnmarshalDiscoveryResponseXml(t *testing.T) {
	response := "<LogicBox swVersion=\"2.5.0\" hwVersion=\"1.0.0\" mac=\"54:10:EC:85:28:BB\" protocol=\"MCP V3.0\"/>"

	expected := Gateway{
		SoftwareVersion: "2.5.0",
		HardwareVersion: "1.0.0",
		MacAddress:      "54:10:EC:85:28:BB",
		Protocol:        "MCP V3.0",
		//IpAddress:       "1.2.3.4",
		//Port:            4001,
	}

	decoded := Gateway{}

	err := xml.Unmarshal([]byte(response), &decoded)
	if err != nil {
		t.Logf("failed to unmarshall received xml content. %v", err)
		t.Fail()
	}

	if !reflect.DeepEqual(expected, decoded) {
		t.Fail()
		fmt.Printf("Actual:\t\t%+v\nExpected:\t%+v\n", decoded, expected)
	}
}
