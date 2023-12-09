package sdk

import (
	"context"
	"encoding/xml"
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestDiscoveryOnRealGateway(t *testing.T) {
	t.SkipNow()

	ctx := context.Background()
	discovery := NewDiscovery(ctx, func(gateway Gateway) {
		t.Logf("Response received: %+v\n", gateway)
	})

	t.Logf("Start discovery...\n")
	discovery.Start()

	t.Logf("Waiting few second for responses...\n")
	time.Sleep(time.Second * 20)

	list := discovery.GetList()
	t.Logf("list: %+v\n", list)

	t.Logf("Stop disovery...\n")
	discovery.Stop()
	t.Logf("Terminated\n")
}

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
		fmt.Errorf("failed to unmarshall received xml content. %v", err)
		t.Fail()
	}

	if !reflect.DeepEqual(expected, decoded) {
		t.Fail()
		fmt.Printf("Actual:\t\t%+v\nExpected:\t%+v\n", decoded, expected)
	}
}
