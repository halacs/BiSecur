package sdk

import (
	"context"
	"encoding/xml"
	"fmt"
	"net"
	"sync"
	"time"
)

const (
	dicoverGatewayPacket = "<Discover target=\"Gateway\" />"
)

// <Discover target="Gateway" />
// <LogicBox swVersion="2.5.0" hwVersion="1.0.0" mac="54:10:EC:85:28:BB" protocol="MCP V3.0"/>
type Gateway struct {
	SoftwareVersion string `xml:"swVersion,attr"`
	HardwareVersion string `xml:"hwVersion,attr"`
	MacAddress      string `xml:"mac,attr"`
	Protocol        string `xml:"protocol,attr"`
	IpAddress       net.IP
	Port            int
}

type Discovery struct {
	wg           *sync.WaitGroup
	ctx          context.Context
	cancelFunc   context.CancelFunc
	callbackFunc func(gateway Gateway)
	gateways     map[string]Gateway
}

func NewDiscovery(ctx context.Context, callbackFunc func(gateway Gateway)) *Discovery {
	ctx2, cancelFunc := context.WithCancel(ctx)
	return &Discovery{
		wg:           &sync.WaitGroup{},
		ctx:          ctx2,
		cancelFunc:   cancelFunc,
		callbackFunc: callbackFunc,
		gateways:     make(map[string]Gateway),
	}
}

func (d *Discovery) sendDiscoveryPacket(broadcastConn *net.UDPConn) {
	_, err := broadcastConn.Write([]byte(dicoverGatewayPacket))
	if err != nil {
		fmt.Errorf("failed to send discovery packet. %v", err)
	}
	fmt.Println("Packet sent")
}

func (d *Discovery) sendingPackets() {
	d.wg.Add(1)
	go func() {
		defer d.wg.Done()

		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		broadcastAddr, err := net.ResolveUDPAddr("udp", "255.255.255.255:4001")
		if err != nil {
			fmt.Errorf("failed to resolve UDP address. %v", err)
			return
		}

		broadcastConn, err := net.DialUDP("udp", nil, broadcastAddr)
		if err != nil {
			fmt.Errorf("failed to open connection. %v", err)
			return
		}

		defer func() {
			err := broadcastConn.Close()
			fmt.Errorf("failed to close connection. %v", err)
		}()

		d.sendDiscoveryPacket(broadcastConn)

		for {
			select {
			case <-d.ctx.Done():
				return
			case <-ticker.C:
				d.sendDiscoveryPacket(broadcastConn)
			}
		}
	}()

}

func (d *Discovery) receivingPackets() {
	d.wg.Add(1)
	go func() {
		defer d.wg.Done()

		connection, err := net.ListenPacket("udp", ":4002")
		if err != nil {
			fmt.Errorf("failed to open listening socket. %v", err)
		}
		defer func() {
			err := connection.Close()
			fmt.Errorf("failed to close network connection. %v", err)
		}()

		for {
			select {
			case <-d.ctx.Done():
				return
			default:
				data := make([]byte, 10000)

				//connection.SetReadDeadline(time.Now().Add(5 * time.Second))
				_, addr, err := connection.ReadFrom(data)
				if err != nil {
					continue
				}
				netUdpAddr := addr.(*net.UDPAddr)

				found := Gateway{
					//SoftwareVersion: "",
					//HardwareVersion: "",
					//MacAddress:      "",
					//Protocol:        "",
					IpAddress: netUdpAddr.IP,
					Port:      netUdpAddr.Port,
				}

				err = xml.Unmarshal(data, &found)
				if err != nil {
					fmt.Errorf("failed to unmarshall received xml content. %v", err)
				}

				key := fmt.Sprintf("%s%d", found.IpAddress, found.Port)
				d.gateways[key] = found

				if d.callbackFunc != nil {
					d.callbackFunc(found)
				}
			}
		}
	}()
}

func (d *Discovery) Start() error {
	d.receivingPackets()
	d.sendingPackets()
	return nil
}

func (d *Discovery) GetList() []Gateway {
	gl := make([]Gateway, 0)
	for _, g := range d.gateways {
		gl = append(gl, g)
	}

	return gl
}

func (d *Discovery) Stop() error {
	d.cancelFunc()
	d.wg.Wait()
	return nil
}
