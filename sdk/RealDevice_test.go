//go:build realdevice

package sdk

import (
	"context"
	"os"
	"strconv"
	"testing"
	"time"
)

var (
	sourceMacAddress      = [6]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	destinationMacAddress = [6]byte{0x54, 0x10, 0xEC, 0x85, 0x28, 0xBB}

	host     = ""
	port     = 0
	username = ""
	password = ""
)

func init() {
	// Using environment variables to avoid leaking sensitive information.
	host = os.Getenv("host")

	var err error
	portStr := os.Getenv("port")
	port, err = strconv.Atoi(portStr)
	if err != nil {
		panic(err)
	}

	username = os.Getenv("username")
	password = os.Getenv("password")
}

func TestGetMacOnRealGateway(t *testing.T) {
	client := NewClient(sourceMacAddress, destinationMacAddress, host, port)
	err := client.Open()
	if err != nil {
		t.Logf("%v", err)
		t.Fail()
	}

	defer func() {
		err2 := client.Close()
		if err2 != nil {
			t.Logf("%v", err2)
			t.Fail()
		}
	}()

	mac, err := client.GetMac()
	if err != nil {
		t.Logf("%v", err)
		t.Fail()
	}

	t.Logf("Received MAC address: %X", mac)
}

func TestGetNameOnRealGateway(t *testing.T) {
	client := NewClient(sourceMacAddress, destinationMacAddress, host, port)
	err := client.Open()
	if err != nil {
		t.Logf("%v", err)
		t.Fail()
	}

	defer func() {
		err2 := client.Close()
		if err2 != nil {
			t.Logf("%v", err2)
			t.Fail()
		}
	}()

	name, err := client.GetName()
	if err != nil {
		t.Logf("%v", err)
		t.Fail()
	}

	t.Logf("Received name: %s", name)
}

func TestGetGroupsOnRealGateway(t *testing.T) {
	client := NewClient(sourceMacAddress, destinationMacAddress, host, port)
	err := client.Open()
	if err != nil {
		t.Logf("%v", err)
		t.Fail()
	}

	defer func() {
		err2 := client.Close()
		if err2 != nil {
			t.Logf("%v", err2)
			t.Fail()
		}
	}()

	groups, err := client.GetGroups()
	if err != nil {
		t.Logf("%v", err)
		t.Fail()
	}

	t.Logf("Received groups: %s", groups.toString())
}

func TestPingOnRealGateway(t *testing.T) {
	client := NewClient(sourceMacAddress, destinationMacAddress, host, port)
	err := client.Open()
	if err != nil {
		t.Logf("%v", err)
		t.Fail()
	}

	defer func() {
		err2 := client.Close()
		if err2 != nil {
			t.Logf("%v", err2)
			t.Fail()
		}
	}()

	err = client.Ping(5)
	if err != nil {
		t.Logf("%v", err)
		t.Fail()
	}
}

func TestLoginOnRealGateway(t *testing.T) {
	client := NewClient(sourceMacAddress, destinationMacAddress, host, port)
	err := client.Open()
	if err != nil {
		t.Logf("%v", err)
		t.Fail()
	}

	defer func() {
		err2 := client.Close()
		if err2 != nil {
			t.Logf("%v", err2)
			t.Fail()
		}
	}()

	err = client.Login(username, password)
	if err != nil {
		t.Logf("%v", err)
		t.Fail()
	}

	t.Logf("client: %+v", client)
}

func TestDiscoveryOnRealGateway(t *testing.T) {
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
