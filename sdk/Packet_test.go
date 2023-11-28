package sdk

import (
	"bisecure/sdk/payload"
	"bytes"
	"encoding/hex"
	"testing"
)

func TestPacketDecode(t *testing.T) {
	testCases := []struct {
		Name                 string
		EncodedInput         string
		ExpectedDecodedInput Packet
	}{
		{
			Name:         "Decode ping 1",
			EncodedInput: "0100000000000A",
			ExpectedDecodedInput: Packet{
				PacketPre: PacketPre{
					TAG:       1,
					Token:     uint32(0x00000000),
					CommandID: COMMANDID_PING,
				},
				payload: payload.EmptyPayload(),
				PacketPost: PacketPost{
					Checksum: 0x0A,
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(test *testing.T) {
			encodedInput := testCase.EncodedInput
			encodedInputBytes, err := hex.DecodeString(encodedInput)

			encodedInputBuffer := new(bytes.Buffer)
			_, err = encodedInputBuffer.Write(encodedInputBytes)
			if err != nil {
				test.Logf("Unexcepted error happend. %v", err)
				test.Fail()
			}

			packetLen := uint16(encodedInputBuffer.Len())

			decoded, err := DecodePacket(packetLen, encodedInputBuffer)
			if err != nil {
				test.Logf("Unexcepted error happend. %v", err)
				test.Fail()
			}

			expected := testCase.ExpectedDecodedInput
			if decoded == nil || !expected.Equal(decoded) {
				test.Logf("Expected value: 0x%X, Actual value: 0x%X", expected, decoded)
				test.Fail()
			}
		})
	}
}

func TestPacketChecksum(t *testing.T) {
	testCases := []struct {
		Name        string
		Payload     Packet
		ExpectedCrc byte
	}{
		{
			Name: "Checksum ping 1",
			Payload: Packet{
				PacketPre: PacketPre{
					TAG:       1,
					Token:     uint32(0x00000000),
					CommandID: COMMANDID_PING,
				},
				payload: payload.EmptyPayload(),
				PacketPost: PacketPost{
					Checksum: 0x00,
				},
			},
			ExpectedCrc: 0x0A,
		},
		{
			Name: "Checksum ping 2",
			Payload: Packet{
				PacketPre: PacketPre{
					TAG:       2,
					Token:     uint32(0x00000000),
					CommandID: COMMANDID_PING,
				},
				payload: payload.EmptyPayload(),
				PacketPost: PacketPost{
					Checksum: 0x00,
				},
			},
			ExpectedCrc: 0x0B,
		},
		{
			Name: "Checksum ping 3",
			Payload: Packet{
				PacketPre: PacketPre{
					TAG:       3,
					Token:     uint32(0x00000000),
					CommandID: COMMANDID_PING,
				},
				payload: payload.EmptyPayload(),
				PacketPost: PacketPost{
					Checksum: 0x00,
				},
			},
			ExpectedCrc: 0x0C,
		},
		{
			Name: "Checksum Login",
			Payload: Packet{
				PacketPre: PacketPre{
					TAG:       0,
					Token:     uint32(0x00000000),
					CommandID: COMMANDID_LOGIN,
				},
				payload: payload.LoginPayload("thomas", "aaabbbccc"),
				PacketPost: PacketPost{
					Checksum: 0x00,
				},
			},
			ExpectedCrc: 0x2D,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(test *testing.T) {
			payload := testCase.Payload
			crc, err := payload.getChecksum()
			if err != nil {
				test.Logf("Unexcepted error happend. %v", err)
				test.Fail()
			}

			expected := testCase.ExpectedCrc
			if crc != expected {
				test.Logf("Expected value: 0x%X, Actual value: 0x%X", expected, crc)
				test.Fail()
			}
		})
	}
}
