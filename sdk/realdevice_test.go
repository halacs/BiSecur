//go:build realdevice

package sdk

import (
	"context"
	"testing"
	"time"
)

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
