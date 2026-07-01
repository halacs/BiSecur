package sdk

import (
	"bytes"
	"fmt"
	"halsecur/sdk/payload"
	"net"
	"time"

	"github.com/sirupsen/logrus"
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
	return c.transmitCommand(requestTc, true, 5*time.Second)
}

func (c *Client) transmitCommandWithResponseTimeout(requestTc *TransmissionContainer, timeout time.Duration) (*TransmissionContainer, error) {
	return c.transmitCommand(requestTc, true, timeout)
}

func (c *Client) transmitCommandWithNoResponse(requestTc *TransmissionContainer) error {
	_, err := c.transmitCommand(requestTc, false, 0)
	return err
}

// receiveResponse reads and decodes a single response from the connection.
func (c *Client) receiveResponse(timeout time.Duration) (*TransmissionContainer, error) {
	receivedBytesTmp := make([]byte, 10240)
	err := c.connection.SetReadDeadline(time.Now().Add(timeout))
	if err != nil {
		return nil, fmt.Errorf("failed to set read deadline. %v", err)
	}
	size, err := c.connection.Read(receivedBytesTmp)
	if err != nil {
		return nil, fmt.Errorf("failed to read network stream. %v", err)
	}

	c.log.Debugf("Length of received bytes: %d", size)
	c.log.Debugf("Response bytes: %X", receivedBytesTmp[0:size])

	buffer := new(bytes.Buffer)
	_, err = buffer.Write(receivedBytesTmp[0:size])
	if err != nil {
		return nil, fmt.Errorf("failed to write into buffer. %v", err)
	}

	receivedTc, err := DecodeTransmissionContainer(buffer)
	if err != nil {
		return nil, fmt.Errorf("failed to decode transmission container. %v", err)
	}

	c.log.Debugf("Received TC: %v", receivedTc)
	return receivedTc, nil
}

func (c *Client) transmitCommand(requestTc *TransmissionContainer, expectResponse bool, readTimeout time.Duration) (*TransmissionContainer, error) {
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

	if !expectResponse {
		return nil, nil
	}

	return c.receiveResponse(readTimeout)
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

	err = responseTc.isResponseFor(requestTc)
	if err != nil {
		return sendTimestamp, receivedTimestamp, fmt.Errorf("received unexpected packet. %v. %v", responseTc, err)
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

	if response.isResponseFor(tc) != nil {
		return deviceMac, fmt.Errorf("received unexpected packet: %s", response)
	}

	getMacResponsePayload, err := castIfNotError[*payload.GetMac](response)
	if err != nil {
		return [6]byte{}, err
	}
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

	if response.isResponseFor(tc) != nil {
		return "", fmt.Errorf("received unexpected packet: %s", response)
	}

	getMacResponsePayload, err := castIfNotError[*payload.GetNameResponse](response)
	if err != nil {
		return "", err
	}
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

	if response.isResponseFor(tc) != nil {
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

	if response.isResponseFor(tc) != nil {
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

	if response.isResponseFor(tc) != nil {
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

	if response.isResponseFor(tc) != nil {
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
		return fmt.Errorf("failed to send login request. %v", err)
	}

	if response == nil {
		return fmt.Errorf("unexpected nil response value")
	}

	// The gateway sends a LOGOUT frame for each stale session before the
	// actual LOGIN response. Drain these before processing the login.
	const maxDrain = 20
	for i := 0; i < maxDrain; i++ {
		if response.Packet.getCommandID() != COMMANDID_LOGOUT {
			break
		}
		c.log.Debugf("Draining stale LOGOUT response (%d)", i+1)
		response, err = c.receiveResponse(5 * time.Second)
		if err != nil {
			return fmt.Errorf("failed to read response while draining stale sessions. %v", err)
		}
	}

	err = response.isResponseFor(tc)
	if err != nil {
		return fmt.Errorf("received unexpected packet (not the response waiting for). %s", response)
	}

	loginResponse, ok := response.Packet.payload.(*payload.LoginResponse)
	if !ok {
		return fmt.Errorf("received unexpected packet (typecast failed): %s", response)
	}

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

	if response.isResponseFor(tc) != nil {
		return fmt.Errorf("received unexpected packet: %s", response)
	}

	// The gateway acknowledges SetState with an ERROR frame (e.g. PERMISSION_DENIED
	// when the session token has expired). Surface it instead of silently succeeding.
	if err := isErrorResponse(response); err != nil {
		return err
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

	transitionResponse, err := castIfNotError[*payload.HmGetTransitionResponse](response)
	if err != nil {
		return nil, err
	}

	if response.isResponseFor(tc) != nil {
		return transitionResponse, fmt.Errorf("received unexpected packet (not the response waiting for): %s", response)
	}

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

	err = response.isResponseFor(tc)
	if err != nil {
		return 0, fmt.Errorf("received unexpected packet: %s. %v", response, err)
	}

	transitionResponse, err := castIfNotError[*payload.AddUserResponse](response)
	if err != nil {
		return 0, err
	}

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

	if response.isResponseFor(tc) != nil {
		return fmt.Errorf("received unexpected packet: %s", response)
	}

	transitionResponse, err := castIfNotError[*payload.RemoveUserResponse](response)
	if err != nil {
		return err
	}

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

	if response.isResponseFor(tc) != nil {
		return fmt.Errorf("received unexpected packet: %s", response)
	}

	return nil
}

// AddPort enters receive mode and listens for a hand remote's radio signal (ADD_PORT 0x29).
// The hand remote must be pressed within ~10cm of the gateway. Timeout: ~40 seconds.
// In testing, the cloned port did not support position feedback via HM_GET_TRANSITION.
func (c *Client) AddPort() (portId byte, err error) {
	tc := c.getTransmissionContainer(COMMANDID_ADD_PORT, payload.EmptyPayload())
	response, err := c.transmitCommandWithResponseTimeout(tc, 45*time.Second)
	if err != nil {
		return 0, fmt.Errorf("failed to encode packet. %v", err)
	}

	if response == nil {
		return 0, fmt.Errorf("unexpected nil response value")
	}

	err = response.isResponseFor(tc)
	if err != nil {
		return 0, fmt.Errorf("received unexpected packet: %s. %v", response, err)
	}

	portIdResponse, err := castIfNotError[*payload.PortIdResponse](response)
	if err != nil {
		return 0, err
	}

	return portIdResponse.GetPortId(), nil
}

// InheritPort transmits the gateway's own radio code for the motor to learn (INHERIT_PORT 0x41).
// The motor must be in learn mode (press P button). The gateway should be within
// ~20-30cm of the motor for pairing. Timeout: ~40 seconds.
// In testing, the inherited port supported position feedback via HM_GET_TRANSITION.
func (c *Client) InheritPort() (portId byte, err error) {
	tc := c.getTransmissionContainer(COMMANDID_INHERIT_PORT, payload.EmptyPayload())
	response, err := c.transmitCommandWithResponseTimeout(tc, 45*time.Second)
	if err != nil {
		return 0, fmt.Errorf("failed to encode packet. %v", err)
	}

	if response == nil {
		return 0, fmt.Errorf("unexpected nil response value")
	}

	err = response.isResponseFor(tc)
	if err != nil {
		return 0, fmt.Errorf("received unexpected packet: %s. %v", response, err)
	}

	portIdResponse, err := castIfNotError[*payload.PortIdResponse](response)
	if err != nil {
		return 0, err
	}

	return portIdResponse.GetPortId(), nil
}

// RemovePort deletes a paired port from the gateway (REMOVE_PORT 0x42).
func (c *Client) RemovePort(portId byte) error {
	tc := c.getTransmissionContainer(COMMANDID_REMOVE_PORT, payload.RemovePortPayload(portId))
	response, err := c.transmitCommandWithResponse(tc)
	if err != nil {
		return fmt.Errorf("failed to encode packet. %v", err)
	}

	if response == nil {
		return fmt.Errorf("unexpected nil response value")
	}

	if response.isResponseFor(tc) != nil {
		return fmt.Errorf("received unexpected packet: %s", response)
	}

	portIdResponse, err := castIfNotError[*payload.PortIdResponse](response)
	if err != nil {
		return err
	}

	if portIdResponse.GetPortId() != portId {
		return fmt.Errorf("failed to remove port. %v", portIdResponse)
	}

	return nil
}

// GetPorts returns the list of port IDs configured on the gateway (GET_PORTS 0x30).
func (c *Client) GetPorts() ([]byte, error) {
	tc := c.getTransmissionContainer(COMMANDID_GET_PORTS, payload.EmptyPayload())
	response, err := c.transmitCommandWithResponse(tc)
	if err != nil {
		return nil, fmt.Errorf("failed to encode packet. %v", err)
	}

	if response == nil {
		return nil, fmt.Errorf("unexpected nil response value")
	}

	if response.isResponseFor(tc) != nil {
		return nil, fmt.Errorf("received unexpected packet: %s", response)
	}

	portsResponse, err := castIfNotError[*payload.GetPortsResponse](response)
	if err != nil {
		return nil, err
	}

	return (*portsResponse).GetPortIds(), nil
}

// GetType returns the type of a port (GET_TYPE 0x31). See PORT_TYPE_* constants.
func (c *Client) GetType(portId byte) (byte, error) {
	tc := c.getTransmissionContainer(COMMANDID_GET_TYPE, payload.GetTypePayload(portId))
	response, err := c.transmitCommandWithResponse(tc)
	if err != nil {
		return 0, fmt.Errorf("failed to encode packet. %v", err)
	}

	if response == nil {
		return 0, fmt.Errorf("unexpected nil response value")
	}

	if response.isResponseFor(tc) != nil {
		return 0, fmt.Errorf("received unexpected packet: %s", response)
	}

	typeResponse, err := castIfNotError[*payload.GetTypeResponse](response)
	if err != nil {
		return 0, err
	}
	return typeResponse.GetPortType(), nil
}

/*
func (c *Client) SetUserRights(userId byte, ???) error {
	return nil
}

func (c *Client) GetUserRights(userId byte, ???) error {
	return nil
}
*/

// Case response to ErrorResponse if possible, otherwise cast response to type T if type match
func castIfNotError[T payload.PayloadInterface](response *TransmissionContainer) (T, error) {
	var t T

	err := isErrorResponse(response)
	if err != nil {
		return t, err
	}

	t, ok := response.Packet.payload.(T)
	if !ok {
		return t, fmt.Errorf("received unexpected packet (typecast failed): %s", response)
	}
	return t, nil
}

func isErrorResponse(response *TransmissionContainer) error {
	errorResponse, isErrorResponseType := response.Packet.payload.(*payload.ErrorResponse)
	if isErrorResponseType {
		// Return the typed *payload.ErrorResponse (it implements error) so callers
		// can inspect the gateway error code, e.g. to react to PERMISSION_DENIED.
		return errorResponse
	}
	return nil
}
