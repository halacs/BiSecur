package sdk

import (
	"context"
	"encoding/xml"
	"fmt"
	"github.com/sirupsen/logrus"
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
	log          *logrus.Logger
	cancelFunc   context.CancelFunc
	callbackFunc func(gateway Gateway)
	gateways     map[string]Gateway
}

func NewDiscovery(ctx context.Context, log *logrus.Logger, callbackFunc func(gateway Gateway)) *Discovery {
	ctx2, cancelFunc := context.WithCancel(ctx)
	return &Discovery{
		wg:           &sync.WaitGroup{},
		ctx:          ctx2,
		log:          log,
		cancelFunc:   cancelFunc,
		callbackFunc: callbackFunc,
		gateways:     make(map[string]Gateway),
	}
}

func (d *Discovery) sendDiscoveryPacket(broadcastConn *net.UDPConn) error {
	_, err := broadcastConn.Write([]byte(dicoverGatewayPacket))
	if err != nil {
		return fmt.Errorf("failed to send discovery packet. %v", err)
	}
	d.log.Debugf("Packet sent")
	return nil
}

func (d *Discovery) sendingPackets() {
	d.wg.Add(1)
	go func() {
		defer d.wg.Done()

		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		broadcastAddr, err := net.ResolveUDPAddr("udp", "255.255.255.255:4001")
		if err != nil {
			d.log.Errorf("failed to resolve UDP address. %v", err)
			return
		}

		broadcastConn, err := net.DialUDP("udp", nil, broadcastAddr)
		if err != nil {
			d.log.Errorf("failed to open connection. %v", err)
			return
		}

		defer func() {
			err := broadcastConn.Close()
			if err != nil {
				d.log.Errorf("failed to close connection. %v", err)
			}
		}()

		err = d.sendDiscoveryPacket(broadcastConn)
		if err != nil {
			d.log.Errorf("failed to send discovery packet. %v", err)
		}

		for {
			select {
			case <-d.ctx.Done():
				return
			case <-ticker.C:
				err := d.sendDiscoveryPacket(broadcastConn)
				if err != nil {
					d.log.Errorf("failed to send discovery packet. %v", err)
				}
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
			d.log.Errorf("failed to open listening socket. %v", err)
		}
		defer func() {
			err := connection.Close()
			if err != nil {
				d.log.Errorf("failed to close network connection. %v", err)
			}
		}()

		for {
			select {
			case <-d.ctx.Done():
				return
			default:
				data := make([]byte, 10000)

				err = connection.SetReadDeadline(time.Now().Add(5 * time.Second))
				if err != nil {
					d.log.Errorf("failed to set read deadline. %v", err)
					continue
				}
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
					d.log.Errorf("failed to unmarshall received xml content. %v", err)
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
