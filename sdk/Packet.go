package sdk

import (
	"bisecur/sdk/payload"
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strings"
)

type Packet struct {
	PacketPre
	payload payload.PayloadInterface // Command/response payload - variable size
	PacketPost
}

type PacketPre struct {
	TAG       byte   // ?
	Token     uint32 // Authentication token - 4 bytes (Session ID)
	CommandID byte   // Command ID - 1 byte
}

type PacketPost struct {
	Checksum byte // Checksum - 1 byte
}

func (p *PacketPre) getSize() int {
	return binary.Size(p)
}

// Gives back if a package contains data of a response or not
// based on the most significant bit of the raw command ID
func (p *PacketPre) isResponse() bool {
	return (p.CommandID & 0x80) == 0x80
}

// Gives back the command ID after clearing
// the most significant bit set in case of a response
func (p *PacketPre) getCommandID() uint8 {
	return p.CommandID & (RESPONSE_MASK ^ 0xFF)
}

func (p *PacketPost) getSize() int {
	return binary.Size(p)
}

func DecodePacket(packetLength uint16, buffer *bytes.Buffer) (*Packet, error) {
	p := Packet{}

	err := binary.Read(buffer, binary.BigEndian, &p.PacketPre)
	if err != nil {
		return nil, fmt.Errorf("failed to decode part of TransmissionContainerPre. %v", err)
	}

	payloadSize := int(packetLength) - (binary.Size(p.PacketPre) + binary.Size(p.PacketPost)) // don't read last checksum byte
	payloadBytes := make([]byte, payloadSize)
	_, err = buffer.Read(payloadBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to decode payload part of Packet. %v", err)
	}

	// Let's decode packet payload according to Command ID
	commandID := p.CommandID
	switch commandID {
	case COMMANDID_LOGIN:
		pl, err2 := payload.DecodeLoginPayload(payloadBytes)
		if err2 != nil {
			return nil, err2
		}
		p.payload = pl
	case COMMANDID_LOGIN_RESPONSE:
		pl, err2 := payload.DecodeLoginResponsePayload(payloadBytes)
		if err2 != nil {
			return nil, err2
		}
		p.payload = pl
	case COMMANDID_GET_MAC_RESPONSE:
		pl, err2 := payload.DecodeGetMacResponsePayload(payloadBytes)
		if err2 != nil {
			return nil, err2
		}
		p.payload = pl
	case COMMANDID_GET_NAME_RESPONSE:
		pl, err2 := payload.DecodeGetNameResponsePayload(payloadBytes)
		if err2 != nil {
			return nil, err2
		}
		p.payload = pl
	case COMMANDID_ERROR:
		pl, err2 := payload.DecodeErrorPayload(payloadBytes)
		if err2 != nil {
			return nil, err2
		}
		p.payload = pl
	case COMMANDID_JMCP_RESPONSE:
		pl, err2 := payload.DecodeJcmpPayload(payloadBytes)
		if err2 != nil {
			return nil, err2
		}
		p.payload = pl
	case COMMANDID_HM_GET_TRANSITION_RESPONSE:
		pl, err2 := payload.DecodeHmGetTransitionResponsePayload(payloadBytes)
		if err2 != nil {
			return nil, err2
		}
		p.payload = pl
	default:
		p.payload = payload.EmptyPayload()
	}

	err = binary.Read(buffer, binary.BigEndian, &p.PacketPost)
	if err != nil {
		return nil, fmt.Errorf("failed to decode part of TransmissionContainerPre. %v", err)
	}

	return &p, nil
}

func (p *Packet) EncodeToHexString() (string, error) {
	byteData, err := p.encodeIntoHexOrNot(true)
	return strings.ToUpper(string(byteData)), err
}

func (p *Packet) Encode() ([]byte, error) {
	byteData, err := p.encodeIntoHexOrNot(false)
	return byteData, err
}

func (p *Packet) toHexString(data []byte) string {
	return strings.ToUpper(hex.EncodeToString(data))
}

func (p *Packet) encodeIntoHexOrNot(encodeIntoHex bool) ([]byte, error) {
	var err error
	p.Checksum, err = p.getChecksum()
	if err != nil {
		return []byte{}, fmt.Errorf("failed to calculage packet checksum. %v", err)
	}

	buffer := new(bytes.Buffer)

	tmpBuff := new(bytes.Buffer)
	err = binary.Write(tmpBuff, binary.BigEndian, p.PacketPre)
	if err != nil {
		return []byte{}, fmt.Errorf("%v", err)
	}
	if encodeIntoHex {
		buffer.Write([]byte(p.toHexString(tmpBuff.Bytes())))
	} else {
		buffer.Write(tmpBuff.Bytes())
	}

	payloadBytes := p.payload.Encode()
	err = binary.Write(buffer, binary.BigEndian, payloadBytes)
	if err != nil {
		return []byte{}, fmt.Errorf("%v", err)
	}

	tmpBuff = new(bytes.Buffer)
	err = binary.Write(tmpBuff, binary.BigEndian, p.PacketPost)
	if err != nil {
		return []byte{}, fmt.Errorf("%v", err)
	}
	if encodeIntoHex {
		buffer.Write([]byte(p.toHexString(tmpBuff.Bytes())))
	} else {
		buffer.Write(tmpBuff.Bytes())
	}

	return buffer.Bytes(), nil
}

func (p *Packet) GetFrameLength() int {
	return p.PacketPre.getSize() + p.PacketPost.getSize()
}

func (p *Packet) GetLength() int {
	frameLength := p.GetFrameLength()
	dataLength := p.payload.Length()
	return frameLength + dataLength
}

func (p *Packet) GetLengthHexString() string {
	payloadLengthBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(payloadLengthBytes, uint16(p.GetLength()))
	return strings.ToUpper(hex.EncodeToString(payloadLengthBytes))
}

func (p *Packet) getChecksum() (byte, error) {
	value := byte(p.GetLength())

	// Summarize the head bytes
	headBuffer := new(bytes.Buffer)
	err := binary.Write(headBuffer, binary.BigEndian, p.PacketPre)
	if err != nil {
		return 0, fmt.Errorf("%v", err)
	}
	for _, v := range headBuffer.Bytes() {
		value = value + v
	}

	// Summarize the payload bytes
	for _, v := range p.payload.ToByteArray() {
		value = value + v
	}

	value = value & 255 // grab only the lowest byte and drop the rest
	return value, nil
}

func (p *Packet) Equal(o *Packet) bool {
	if p.PacketPre != o.PacketPre {
		return false
	}

	if p.PacketPost != o.PacketPost {
		return false
	}

	if !bytes.Equal(p.payload.ToByteArray(), o.payload.ToByteArray()) {
		return false
	}

	return true
}

func (p *Packet) String() string {
	return fmt.Sprintf("Tag=0x%X, Token=0x%X, CommandID=0x%X (0x%X), payload=[%s], Checksum=0x%X, isResponse=%t",
		p.TAG, p.Token, p.getCommandID(), p.CommandID, p.payload, p.Checksum, p.isResponse())
}
