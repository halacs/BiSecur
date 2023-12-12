package sdk

import (
	"bisecure/sdk/payload"
	"bytes"
	"fmt"
	"net"
	"time"
)

type Client struct {
	destinationMacAddress [6]byte
	sourceMacAddress      [6]byte
	host                  string
	port                  int
	tag                   byte
	token                 uint32
	senderID              byte
	connection            *net.TCPConn
}

func NewClient(sourceMacAddress [6]byte, destinationMacAddress [6]byte, host string, port int) *Client {
	return &Client{
		sourceMacAddress:      sourceMacAddress,
		destinationMacAddress: destinationMacAddress,
		host:                  host,
		port:                  port,
		tag:                   1,
		token:                 0,
		senderID:              0,
	}
}

func (c *Client) getTransmissionContainer(commandID byte, payload payload.PayloadInterface) *TransmissionContainer {
	tag := c.tag
	c.tag = c.tag + 1

	tc := TransmissionContainer{
		TransmissionContainerPre: TransmissionContainerPre{
			SrcMac: c.sourceMacAddress,
			DstMac: c.destinationMacAddress,
		},
		Packet: Packet{
			PacketPre: PacketPre{
				TAG:       tag,
				Token:     c.token,
				CommandID: commandID,
			},
			payload:    payload,
			PacketPost: PacketPost{},
		},
		TransmissionContainerPost: TransmissionContainerPost{},
	}

	return &tc
}

func (c *Client) transmitCommand(requestTc *TransmissionContainer) (*TransmissionContainer, error) {
	requestBytes, err := requestTc.Encode()
	if err != nil {
		return nil, fmt.Errorf("failed to encode packet. %v", err)
	}
	fmt.Printf("Request bytes: %X\n", requestBytes)
	_, err = c.connection.Write(requestBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to write into network stream. %v", err)
	}

	receivedBytesTmp := make([]byte, 10240)
	c.connection.SetReadDeadline(time.Now().Add(time.Second * 5))
	size, err := c.connection.Read(receivedBytesTmp)
	receivedHexString := string(receivedBytesTmp[0:size])
	fmt.Printf("Received bytes: %d\nResponse bytes: %s\n", size, receivedHexString)
	if err != nil {
		return nil, fmt.Errorf("failed to read network stream. %v", err)
	}

	buffer := new(bytes.Buffer)
	_, err = buffer.Write(receivedBytesTmp[0:size])
	if err != nil {
		return nil, fmt.Errorf("failed to write into buffer. %v", err)
	}

	receivedTc, err := DecodeTransmissionContainer(buffer)
	if err != nil {
		return nil, fmt.Errorf("failed to decode transmission container. %v", err)
	}

	fmt.Printf("%v\n", receivedTc)

	return receivedTc, err
}

func (c *Client) Open() error {
	if len(c.host) == 0 {
		return fmt.Errorf("'host' value cannot be empty")
	}

	servAddr := fmt.Sprintf("%s:%d", c.host, c.port)
	tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
	if err != nil {
		return fmt.Errorf("resolveTCPAddr failed. %v", err)
	}

	fmt.Printf("Connecting to %s\n", servAddr)

	c.connection, err = net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		c.connection = nil
		return fmt.Errorf("dial failed. %v", err)
	}

	return nil
}

func (c *Client) Close() error {
	if c.connection == nil {
		return fmt.Errorf("connection is closed")
	}

	err := c.connection.Close()
	c.connection = nil
	return err
}

func (c *Client) isOpened() bool {
	return c.connection != nil
}

func (c *Client) Ping(count int) error {
	received := 0

	for i := 0; i < count; i++ {
		requestTc := c.getTransmissionContainer(COMMANDID_PING, payload.EmptyPayload())
		fmt.Printf("requestTC: %v\n", requestTc)
		responseTc, err := c.transmitCommand(requestTc)
		if err != nil {
			return fmt.Errorf("%v", err)
		}

		if responseTc == nil {
			return fmt.Errorf("unexpected nil responseTc value")
		}

		if !responseTc.isResponseFor(requestTc) {
			return fmt.Errorf("received unexpected packet: %s", responseTc)
		}

		received = received + 1
		fmt.Printf("received %d packet(s)\n", received)

		if i < count {
			time.Sleep(time.Second * 2)
			fmt.Println("-------------------")
		}
	}

	if count != received {
		return fmt.Errorf("lost packets. Sent: %d, Received: %d, Ration %f", count, received, (float32(received)/float32(count))*100.0)
	}

	return nil
}

func (c *Client) GetMac() ([6]byte, error) {
	/*
		"Side Note: GET_MAC is used as a keepalive, every 30 seconds the device receives a GET_MAC message."
		Source: https://sec-consult.com/blog/detail/hoermann-opening-doors-for-everyone/
	*/

	deviceMac := [6]byte{0, 0, 0, 0, 0, 0}

	tc := c.getTransmissionContainer(COMMANDID_GET_MAC, payload.EmptyPayload())
	response, err := c.transmitCommand(tc)
	if err != nil {
		return deviceMac, fmt.Errorf("failed to encode packet. %v", err)
	}

	if response == nil {
		return deviceMac, fmt.Errorf("unexpected nil response value")
	}

	if !response.isResponseFor(tc) {
		return deviceMac, fmt.Errorf("received unexpected packet: %s", response)
	}

	getMacResponsePayload := response.Packet.payload.(*payload.GetMac)
	deviceMac = getMacResponsePayload.GetMac()

	return deviceMac, nil
}

func (c *Client) GetName() (string, error) {
	tc := c.getTransmissionContainer(COMMANDID_GET_NAME, payload.EmptyPayload())
	response, err := c.transmitCommand(tc)
	if err != nil {
		return "", fmt.Errorf("failed to encode packet. %v", err)
	}

	if response == nil {
		return "", fmt.Errorf("unexpected nil response value")
	}

	if !response.isResponseFor(tc) {
		return "", fmt.Errorf("received unexpected packet: %s", response)
	}

	getMacResponsePayload := response.Packet.payload.(*payload.GetNameResponse)
	name := getMacResponsePayload.GetName()
	return name, nil
}

func (c *Client) Login(username string, password string) error {
	if len(username) == 0 {
		return fmt.Errorf("'username' value cannot be empty")
	}

	if len(password) == 0 {
		return fmt.Errorf("'password' value cannot be empty")
	}

	tc := c.getTransmissionContainer(COMMANDID_LOGIN, payload.LoginPayload(username, password))
	response, err := c.transmitCommand(tc)
	if err != nil {
		return fmt.Errorf("failed to encode packet. %v", err)
	}

	if response == nil {
		return fmt.Errorf("unexpected nil response value")
	}

	if !response.isResponseFor(tc) {
		return fmt.Errorf("received unexpected packet: %s", response)
	}

	loginResponse := response.Packet.payload.(*payload.LoginResponse)
	c.token = loginResponse.GetToken()
	c.senderID = loginResponse.GetSenderID()

	return nil
}

func (c *Client) Logout() error {
	return nil
}

func (c *Client) SetState() error {
	return nil
}

func (c *Client) GetTransition() error {
	return nil
}
