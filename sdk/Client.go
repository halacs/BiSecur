package sdk

import (
	"bisecur/sdk/payload"
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
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
	log                   *logrus.Logger
}

func NewClient(log *logrus.Logger, sourceMacAddress [6]byte, destinationMacAddress [6]byte, host string, port int, token uint32) *Client {
	return &Client{
		log:                   log,
		sourceMacAddress:      sourceMacAddress,
		destinationMacAddress: destinationMacAddress,
		host:                  host,
		port:                  port,
		tag:                   1,
		token:                 token,
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

func (c *Client) transmitCommandWithResponse(requestTc *TransmissionContainer) (*TransmissionContainer, error) {
	return c.transmitCommand(requestTc, true)
}

func (c *Client) transmitCommandWithNoResponse(requestTc *TransmissionContainer) error {
	_, err := c.transmitCommand(requestTc, false)
	return err
}

func (c *Client) transmitCommand(requestTc *TransmissionContainer, expectResponse bool) (*TransmissionContainer, error) {
	c.log.Debugf("Request: %s", requestTc.String())
	requestBytes, err := requestTc.Encode()
	if err != nil {
		return nil, fmt.Errorf("failed to encode packet. %v", err)
	}
	c.log.Debugf("Request bytes: %X", requestBytes)
	_, err = c.connection.Write(requestBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to write into network stream. %v", err)
	}

	var receivedTc *TransmissionContainer

	if expectResponse {
		receivedBytesTmp := make([]byte, 10240)
		err := c.connection.SetReadDeadline(time.Now().Add(time.Second * 5))
		if err != nil {
			return nil, fmt.Errorf("failed to set read deadline. %v", err)
		}
		size, err := c.connection.Read(receivedBytesTmp)
		receivedHexString := string(receivedBytesTmp[0:size])
		c.log.Debugf("Length of received bytes: %d", size)
		c.log.Debugf("Response bytes: %s", receivedHexString)
		if err != nil {
			return nil, fmt.Errorf("failed to read network stream. %v", err)
		}

		buffer := new(bytes.Buffer)
		_, err = buffer.Write(receivedBytesTmp[0:size])
		if err != nil {
			return nil, fmt.Errorf("failed to write into buffer. %v", err)
		}

		receivedTc, err = DecodeTransmissionContainer(buffer)
		if err != nil {
			return nil, fmt.Errorf("failed to decode transmission container. %v", err)
		}

		c.log.Debugf("Received TC: %v", receivedTc)
	}
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

	c.log.Debugf("Connecting to %s", servAddr)

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

func (c *Client) IsOpened() bool {
	return c.connection != nil
}

func (c *Client) Ping() (int64, int64, error) {
	requestTc := c.getTransmissionContainer(COMMANDID_PING, payload.EmptyPayload())
	c.log.Debugf("requestTC: %v", requestTc)

	sendTimestamp := time.Now().UnixMilli()
	responseTc, err := c.transmitCommandWithResponse(requestTc)
	receivedTimestamp := time.Now().UnixMilli()

	c.log.Debugf("responseTC: %v", responseTc)

	if err != nil {
		return sendTimestamp, receivedTimestamp, fmt.Errorf("%v", err)
	}

	if responseTc == nil {
		return sendTimestamp, receivedTimestamp, fmt.Errorf("unexpected nil responseTc value")
	}

	if !responseTc.isResponseFor(requestTc) {
		return sendTimestamp, receivedTimestamp, fmt.Errorf("received unexpected packet. %v", responseTc)
	}

	return sendTimestamp, receivedTimestamp, nil
}

func (c *Client) GetMac() ([6]byte, error) {
	/*
		"Side Note: GET_MAC is used as a keepalive, every 30 seconds the device receives a GET_MAC message."
		Source: https://sec-consult.com/blog/detail/hoermann-opening-doors-for-everyone/
	*/

	deviceMac := [6]byte{0, 0, 0, 0, 0, 0}

	tc := c.getTransmissionContainer(COMMANDID_GET_MAC, payload.EmptyPayload())
	response, err := c.transmitCommandWithResponse(tc)
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
	response, err := c.transmitCommandWithResponse(tc)
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

// GetGroups returns the Groups are the paired devices. This call returns all devices known to the gateway.
func (c *Client) GetGroups() (*Groups, error) {
	tc := c.getTransmissionContainer(COMMANDID_JMCP, payload.JcmpPayload("{\"CMD\":\"GET_GROUPS\"}"))
	response, err := c.transmitCommandWithResponse(tc)
	if err != nil {
		return nil, fmt.Errorf("failed to encode packet. %v", err)
	}

	if response == nil {
		return nil, fmt.Errorf("unexpected nil response value")
	}

	if !response.isResponseFor(tc) {
		return nil, fmt.Errorf("received unexpected packet: %s", response)
	}

	responsePayload := string(response.Packet.payload.ToByteArray())
	groups, err := DecodeGroups(responsePayload)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal Groups object. %v", err)
	}

	return &groups, nil
}

// GetGroupsForUser returns only the devices that are paired with the current user. We probably should always use this.
func (c *Client) GetGroupsForUser(userID byte) (*Groups, error) {
	tc := c.getTransmissionContainer(COMMANDID_JMCP, payload.JcmpPayload(fmt.Sprintf("{\"CMD\":\"GET_GROUPS\", \"FORUSER\":%d}", userID)))
	response, err := c.transmitCommandWithResponse(tc)
	if err != nil {
		return nil, fmt.Errorf("failed to encode packet. %v", err)
	}

	if response == nil {
		return nil, fmt.Errorf("unexpected nil response value")
	}

	if !response.isResponseFor(tc) {
		return nil, fmt.Errorf("received unexpected packet: %s", response)
	}

	responsePayload := string(response.Packet.payload.ToByteArray())
	groups, err := DecodeGroups(responsePayload)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal Groups object. %v", err)
	}

	return &groups, nil
}

func (c *Client) GetUsers() (*Users, error) {
	tc := c.getTransmissionContainer(COMMANDID_JMCP, payload.JcmpPayload("{\"CMD\":\"GET_USERS\"}"))
	response, err := c.transmitCommandWithResponse(tc)
	if err != nil {
		return nil, fmt.Errorf("failed to encode packet. %v", err)
	}

	if response == nil {
		return nil, fmt.Errorf("unexpected nil response value")
	}

	if !response.isResponseFor(tc) {
		return nil, fmt.Errorf("received unexpected packet: %s", response)
	}

	responsePayload := string(response.Packet.payload.ToByteArray())
	users, err := DecodeUsers(responsePayload)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal Groups object. %v", err)
	}

	return &users, nil
}

// GetValues returns a map of port and some kind of number. I don't know how to handle that.
func (c *Client) GetValues() (*Values, error) {
	tc := c.getTransmissionContainer(COMMANDID_JMCP, payload.JcmpPayload("{\"CMD\":\"GET_VALUES\"}"))
	response, err := c.transmitCommandWithResponse(tc)
	if err != nil {
		return nil, fmt.Errorf("failed to encode packet. %v", err)
	}

	if response == nil {
		return nil, fmt.Errorf("unexpected nil response value")
	}

	if !response.isResponseFor(tc) {
		return nil, fmt.Errorf("received unexpected packet: %s", response)
	}

	responsePayload := string(response.Packet.payload.ToByteArray())
	values, err := DecodeValues(responsePayload)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal Groups object. %v", err)
	}

	return &values, nil
}

func (c *Client) Login(username string, password string) error {
	if len(username) == 0 {
		return fmt.Errorf("'username' value cannot be empty")
	}

	if len(password) == 0 {
		return fmt.Errorf("'password' value cannot be empty")
	}

	tc := c.getTransmissionContainer(COMMANDID_LOGIN, payload.LoginPayload(username, password))
	response, err := c.transmitCommandWithResponse(tc)
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

func (c *Client) SetToken(token uint32) {
	c.token = token
}

func (c *Client) GetToken() uint32 {
	return c.token
}

func (c *Client) Logout() error {
	tc := c.getTransmissionContainer(COMMANDID_LOGOUT, payload.EmptyPayload())
	err := c.transmitCommandWithNoResponse(tc) // Don't care about response, if any. It seems gateway doesn't send response for logout request
	if err != nil {
		return fmt.Errorf("failed to encode packet. %v", err)
	}

	// clear local token store
	c.token = 0

	return nil
}

func (c *Client) SetState(portID byte) error {
	tc := c.getTransmissionContainer(COMMANDID_SET_STATE, payload.SetStatePayload(portID))
	response, err := c.transmitCommandWithResponse(tc)
	if err != nil {
		return fmt.Errorf("failed to encode packet. %v", err)
	}

	if response == nil {
		return fmt.Errorf("unexpected nil response value")
	}

	if !response.isResponseFor(tc) {
		return fmt.Errorf("received unexpected packet: %s", response)
	}

	c.log.Debugf("Set State response: %s", response.String())

	return nil
}

// GetTransition returns the current state of the port. You can see how much open it is or if it is still running.
func (c *Client) GetTransition(portID byte) (*payload.HmGetTransitionResponse, error) {
	tc := c.getTransmissionContainer(COMMANDID_HM_GET_TRANSITION, payload.HmGetTransitionPayload(portID))
	response, err := c.transmitCommandWithResponse(tc)
	if err != nil {
		return nil, fmt.Errorf("failed to encode packet. %v", err)
	}

	if response == nil {
		return nil, fmt.Errorf("unexpected nil response value")
	}

	if !response.isResponseFor(tc) {
		return nil, fmt.Errorf("received unexpected packet: %s", response)
	}

	transitionResponse := response.Packet.payload.(*payload.HmGetTransitionResponse)
	return transitionResponse, nil
}

func (c *Client) AddUser(userName string, password string) (userId byte, err error) {
	tc := c.getTransmissionContainer(COMMANDID_ADD_USER, payload.LoginPayload(userName, password))
	response, err := c.transmitCommandWithResponse(tc)
	if err != nil {
		return 0, fmt.Errorf("failed to encode packet. %v", err)
	}

	if response == nil {
		return 0, fmt.Errorf("unexpected nil response value")
	}

	if !response.isResponseFor(tc) {
		return 0, fmt.Errorf("received unexpected packet: %s", response)
	}

	transitionResponse := response.Packet.payload.(*payload.AddUserResponse)
	newUserId := transitionResponse.GetUserId()
	return newUserId, nil
}

func (c *Client) RemoveUser(userId byte) error {
	tc := c.getTransmissionContainer(COMMANDID_REMOVE_USER, payload.RemoveUserPayload(userId))
	response, err := c.transmitCommandWithResponse(tc)
	if err != nil {
		return fmt.Errorf("failed to encode packet. %v", err)
	}

	if response == nil {
		return fmt.Errorf("unexpected nil response value")
	}

	if !response.isResponseFor(tc) {
		return fmt.Errorf("received unexpected packet: %s", response)
	}

	transitionResponse := response.Packet.payload.(*payload.RemoveUserResponse)
	if transitionResponse.GetUserId() != userId {
		return fmt.Errorf("failed to remove user. %v", transitionResponse)
	}

	return nil
}

func (c *Client) PasswordChange(userId byte, newPassword string) error {
	tc := c.getTransmissionContainer(COMMANDID_CHANGE_PASSWD, payload.ChangeUserPasswordPayload(userId, newPassword))
	response, err := c.transmitCommandWithResponse(tc)
	if err != nil {
		return fmt.Errorf("failed to encode packet. %v", err)
	}

	if response == nil {
		return fmt.Errorf("unexpected nil response value")
	}

	if !response.isResponseFor(tc) {
		return fmt.Errorf("received unexpected packet: %s", response)
	}

	return nil
}

/*
func (c *Client) SetUserRights(userId byte, ???) error {
	return nil
}

func (c *Client) GetUserRights(userId byte, ???) error {
	return nil
}
*/
