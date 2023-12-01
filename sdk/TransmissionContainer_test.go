package sdk

import (
	"bisecure/sdk/payload"
	"bytes"
	"encoding/hex"
	"strings"
	"testing"
)

func TestTransmissionContainerEncode(t *testing.T) {
	testCases := []struct {
		Name                  string
		Request               TransmissionContainer
		ExpectedServerRequest string
	}{
		{
			Name: "Ping 1",
			Request: TransmissionContainer{
				TransmissionContainerPre: TransmissionContainerPre{
					SrcMac: [6]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
					DstMac: [6]byte{0x54, 0x10, 0xEC, 0x85, 0x28, 0xBB},
				},
				Packet: Packet{
					PacketPre: PacketPre{
						TAG:       1,
						Token:     uint32(0x00000000),
						CommandID: COMMANDID_PING,
					},
					payload:    payload.EmptyPayload(),
					PacketPost: PacketPost{
						//Checksum: 0x0A,
					},
				},
				TransmissionContainerPost: TransmissionContainerPost{
					//Checksum: 0x68,
				},
			},
			ExpectedServerRequest: "0000000000005410EC8528BB00090100000000000A68",
		},
		{
			Name: "Ping 2",
			Request: TransmissionContainer{
				TransmissionContainerPre: TransmissionContainerPre{
					SrcMac: [6]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
					DstMac: [6]byte{0x54, 0x10, 0xEC, 0x85, 0x28, 0xBB},
				},
				Packet: Packet{
					PacketPre: PacketPre{
						TAG:       2,
						Token:     uint32(0x00000000),
						CommandID: COMMANDID_PING,
					},
					payload:    payload.EmptyPayload(),
					PacketPost: PacketPost{
						//Checksum: 0x0B,
					},
				},
				TransmissionContainerPost: TransmissionContainerPost{
					//Checksum: 0x6A,
				},
			},
			ExpectedServerRequest: "0000000000005410EC8528BB00090200000000000B6A",
		},
		{
			Name: "Ping 3",
			Request: TransmissionContainer{
				TransmissionContainerPre: TransmissionContainerPre{
					SrcMac: [6]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
					DstMac: [6]byte{0x54, 0x10, 0xEC, 0x85, 0x28, 0xBB},
				},
				Packet: Packet{
					PacketPre: PacketPre{
						TAG:       3,
						Token:     uint32(0x00000000),
						CommandID: COMMANDID_PING,
					},
					payload:    payload.EmptyPayload(),
					PacketPost: PacketPost{
						//Checksum: 0x0C,
					},
				},
				TransmissionContainerPost: TransmissionContainerPost{
					//Checksum: 0x6C,
				},
			},
			ExpectedServerRequest: "0000000000005410EC8528BB00090300000000000C6C",
		},
		{
			Name: "LoginPayload",
			Request: TransmissionContainer{
				TransmissionContainerPre: TransmissionContainerPre{
					SrcMac: [6]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
					DstMac: [6]byte{0x54, 0x10, 0xEC, 0x85, 0x28, 0xBB},
				},
				Packet: Packet{
					PacketPre: PacketPre{
						TAG:       0x00,
						Token:     uint32(0x00000000),
						CommandID: COMMANDID_LOGIN,
					},
					payload:    payload.LoginPayload("username", "password"),
					PacketPost: PacketPost{
						//Checksum: 0x06,
					},
				},
				TransmissionContainerPost: TransmissionContainerPost{
					//Checksum: 0x9F,
				},
			},
			ExpectedServerRequest: "0000000000005410EC8528BB001A00000000001008757365726E616D6570617373776F7264059D",
		},
		{
			Name: "GetNamePayload",
			Request: TransmissionContainer{
				TransmissionContainerPre: TransmissionContainerPre{
					SrcMac: [6]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
					DstMac: [6]byte{0x54, 0x10, 0xEC, 0x85, 0x28, 0xBB},
				},
				Packet: Packet{
					PacketPre: PacketPre{
						TAG:       0x01,
						Token:     uint32(0x00000000),
						CommandID: COMMANDID_GET_NAME,
					},
					payload:    payload.GetNamePayload(),
					PacketPost: PacketPost{
						//Checksum: 0x5E,
					},
				},
				TransmissionContainerPost: TransmissionContainerPost{
					//Checksum: 0xC3,
				},
			},
			ExpectedServerRequest: "0000000000005410EC8528BB00090100000000263062",
		},
		{
			Name: "HmGetTransitionPayload",
			Request: TransmissionContainer{
				TransmissionContainerPre: TransmissionContainerPre{
					SrcMac: [6]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x06},
					DstMac: [6]byte{0x54, 0x10, 0xEC, 0x85, 0x28, 0xBB},
				},
				Packet: Packet{
					PacketPre: PacketPre{
						TAG:       0x06,
						Token:     uint32(0x96833386),
						CommandID: COMMANDID_HM_GET_TRANSITION,
					},
					payload:    payload.HmGetTransitionPayload(0x00),
					PacketPost: PacketPost{
						//Checksum: 0x52,
					},
				},
				TransmissionContainerPost: TransmissionContainerPost{
					//Checksum: 0x06,
				},
			},
			ExpectedServerRequest: "0000000000065410EC8528BB000A069683338670005206",
		},
		{
			Name: "JcmpGetValuesPayload",
			Request: TransmissionContainer{
				TransmissionContainerPre: TransmissionContainerPre{
					SrcMac: [6]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x06},
					DstMac: [6]byte{0x54, 0x10, 0xEC, 0x85, 0x28, 0xBB},
				},
				Packet: Packet{
					PacketPre: PacketPre{
						TAG:       0x04,
						Token:     uint32(0x96833386),
						CommandID: COMMANDID_JMCP,
					},
					payload:    payload.JcmpPayload("{\"cmd\":\"GET_VALUES\"}"),
					PacketPost: PacketPost{
						//Checksum: 0xF6,
					},
				},
				TransmissionContainerPost: TransmissionContainerPost{
					//Checksum: 0x3B,
				},
			},
			ExpectedServerRequest: "0000000000065410EC8528BB001D0496833386067B22636D64223A224745545F56414C554553227DF63B",
		},
		{
			Name: "JcmpGetGroupsPayload",
			Request: TransmissionContainer{
				TransmissionContainerPre: TransmissionContainerPre{
					SrcMac: [6]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x06},
					DstMac: [6]byte{0x54, 0x10, 0xEC, 0x85, 0x28, 0xBB},
				},
				Packet: Packet{
					PacketPre: PacketPre{
						TAG:       0x05,
						Token:     uint32(0x96833386),
						CommandID: COMMANDID_JMCP,
					},
					payload:    payload.JcmpPayload("{\"CMD\":\"GET_GROUPS\"}"),
					PacketPost: PacketPost{
						//Checksum: 0xA7,
					},
				},
				TransmissionContainerPost: TransmissionContainerPost{
					//Checksum: 0x33,
				},
			},
			ExpectedServerRequest: "0000000000065410EC8528BB001D0596833386067B22434D44223A224745545F47524F555053227DA733",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(test *testing.T) {
			commandRequest := testCase.Request
			raw, err := commandRequest.Encode()
			if err != nil {
				test.Logf("Failed to encode request. %v", err)
				test.Fail()
				return
			}

			str := strings.ToUpper(string(raw))
			expected := strings.ToUpper(testCase.ExpectedServerRequest)
			if str != expected {
				test.Logf("Expected value: %v, Actual value: %v", expected, str)
				test.Fail()
			}
		})
	}
}

func TestTransmissionContainerDecode(t *testing.T) {
	testCases := []struct {
		Name                 string
		EncodedInput         string
		ExpectedDecodedInput TransmissionContainer
	}{
		{
			Name:         "Ping 1",
			EncodedInput: "0000000000005410EC8528BB00090100000000000A68",
			ExpectedDecodedInput: TransmissionContainer{
				TransmissionContainerPre: TransmissionContainerPre{
					SrcMac:     [6]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
					DstMac:     [6]byte{0x54, 0x10, 0xEC, 0x85, 0x28, 0xBB},
					BodyLength: 9,
				},
				Packet: Packet{
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
				TransmissionContainerPost: TransmissionContainerPost{
					Checksum: 0x68,
				},
			},
		},
		{
			Name:         "Ping 2",
			EncodedInput: "0000000000005410EC8528BB00090200000000000B6A",
			ExpectedDecodedInput: TransmissionContainer{
				TransmissionContainerPre: TransmissionContainerPre{
					SrcMac:     [6]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
					DstMac:     [6]byte{0x54, 0x10, 0xEC, 0x85, 0x28, 0xBB},
					BodyLength: 9,
				},
				Packet: Packet{
					PacketPre: PacketPre{
						TAG:       2,
						Token:     uint32(0x00000000),
						CommandID: COMMANDID_PING,
					},
					payload: payload.EmptyPayload(),
					PacketPost: PacketPost{
						Checksum: 0x0B,
					},
				},
				TransmissionContainerPost: TransmissionContainerPost{
					Checksum: 0x6A,
				},
			},
		},
		{
			Name:         "Ping 3",
			EncodedInput: "0000000000005410EC8528BB00090300000000000C6C",
			ExpectedDecodedInput: TransmissionContainer{
				TransmissionContainerPre: TransmissionContainerPre{
					SrcMac:     [6]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
					DstMac:     [6]byte{0x54, 0x10, 0xEC, 0x85, 0x28, 0xBB},
					BodyLength: 9,
				},
				Packet: Packet{
					PacketPre: PacketPre{
						TAG:       3,
						Token:     uint32(0x00000000),
						CommandID: COMMANDID_PING,
					},
					payload: payload.EmptyPayload(),
					PacketPost: PacketPost{
						Checksum: 0x0C,
					},
				},
				TransmissionContainerPost: TransmissionContainerPost{
					Checksum: 0x6C,
				},
			},
		},
		{
			Name:         "LoginPayload",
			EncodedInput: "0000000000005410EC8528BB001A00000000001008757365726E616D6570617373776F7264059D",
			ExpectedDecodedInput: TransmissionContainer{
				TransmissionContainerPre: TransmissionContainerPre{
					SrcMac:     [6]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
					DstMac:     [6]byte{0x54, 0x10, 0xEC, 0x85, 0x28, 0xBB},
					BodyLength: 26,
				},
				Packet: Packet{
					PacketPre: PacketPre{
						TAG:       0x00,
						Token:     uint32(0x00000000),
						CommandID: COMMANDID_LOGIN,
					},
					payload: payload.LoginPayload("username", "password"),
					PacketPost: PacketPost{
						Checksum: 0x05,
					},
				},
				TransmissionContainerPost: TransmissionContainerPost{
					Checksum: 0x9D,
				},
			},
		},
		{
			Name:         "LoginFailedResponse",
			EncodedInput: "5410EC8528BB000000000004000A020000000001020FDD",
			ExpectedDecodedInput: TransmissionContainer{
				TransmissionContainerPre: TransmissionContainerPre{
					SrcMac:     [6]byte{0x54, 0x10, 0xEC, 0x85, 0x28, 0xBB},
					DstMac:     [6]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x04},
					BodyLength: 0x0A,
				},
				Packet: Packet{
					PacketPre: PacketPre{
						TAG:       0x02,
						Token:     uint32(0x00000000),
						CommandID: COMMANDID_ERROR,
					},
					payload: payload.ErrorPayload(payload.ERROR_LOGIN_FAILED),
					PacketPost: PacketPost{
						Checksum: 0x0F,
					},
				},
				TransmissionContainerPost: TransmissionContainerPost{
					Checksum: 0xDD,
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(test *testing.T) {
			encodedInput := testCase.EncodedInput
			encodedInputBytes, err := hex.DecodeString(encodedInput)
			if err != nil {
				test.Logf("Failed to decode test data. %v", err)
				test.Fail()
				return
			}

			encodedInputBuffer := new(bytes.Buffer)
			_, err = encodedInputBuffer.Write(encodedInputBytes)
			if err != nil {
				test.Logf("Failed to write test data into buffer. %v", err)
				test.Fail()
				return
			}

			decoded, err := DecodeTransmissionContainer(encodedInputBuffer)
			if err != nil {
				test.Logf("Failed to encode request. %v", err)
				test.Fail()
				return
			}

			expected := &testCase.ExpectedDecodedInput
			if decoded == nil || !expected.Equal(decoded) {
				test.Logf("Expected value: 0x%X, Actual value: 0x%X", expected, decoded)
				test.Fail()
			}
			test.Logf("Decoded transmission container: %v", decoded)
		})
	}
}

func TestTransmissionContainerChecksum(t *testing.T) {
	testCases := []struct {
		Name             string
		Packet           TransmissionContainer
		ExpectedChecksum byte
	}{
		{
			Name: "Ping 1 checksum",
			Packet: TransmissionContainer{
				TransmissionContainerPre: TransmissionContainerPre{
					SrcMac: [6]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
					DstMac: [6]byte{0x54, 0x10, 0xEC, 0x85, 0x28, 0xBB},
				},
				Packet: Packet{
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
				TransmissionContainerPost: TransmissionContainerPost{
					Checksum: 0x00,
				},
			},
			ExpectedChecksum: 0x68,
		},
		{
			Name: "Ping 2 checksum",
			Packet: TransmissionContainer{
				TransmissionContainerPre: TransmissionContainerPre{
					SrcMac: [6]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
					DstMac: [6]byte{0x54, 0x10, 0xEC, 0x85, 0x28, 0xBB},
				},
				Packet: Packet{
					PacketPre: PacketPre{
						TAG:       2,
						Token:     uint32(0x00000000),
						CommandID: COMMANDID_PING,
					},
					payload: payload.EmptyPayload(),
					PacketPost: PacketPost{ // -118
						Checksum: 0x0B,
					},
				},
				TransmissionContainerPost: TransmissionContainerPost{
					Checksum: 0x00,
				},
			},
			ExpectedChecksum: 0x6A,
		},
		{
			Name: "Ping 3 checksum",
			Packet: TransmissionContainer{
				TransmissionContainerPre: TransmissionContainerPre{
					SrcMac: [6]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
					DstMac: [6]byte{0x54, 0x10, 0xEC, 0x85, 0x28, 0xBB},
				},
				Packet: Packet{
					PacketPre: PacketPre{
						TAG:       3,
						Token:     uint32(0x00000000),
						CommandID: COMMANDID_PING,
					},
					payload: payload.EmptyPayload(),
					PacketPost: PacketPost{
						Checksum: 0x0C,
					},
				},
				TransmissionContainerPost: TransmissionContainerPost{
					Checksum: 0x00,
				},
			},
			ExpectedChecksum: 0x6C,
		},
		{
			Name: "LoginPayload checksum",
			Packet: TransmissionContainer{
				TransmissionContainerPre: TransmissionContainerPre{
					SrcMac: [6]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
					DstMac: [6]byte{0x54, 0x10, 0xEC, 0x85, 0x28, 0xBB},
				},
				Packet: Packet{
					PacketPre: PacketPre{
						TAG:       1,
						Token:     uint32(0x00000000),
						CommandID: COMMANDID_LOGIN,
					},
					payload: payload.LoginPayload("thomas", "aaabbbccc"),
					PacketPost: PacketPost{
						Checksum: 0x06,
					},
				},
				TransmissionContainerPost: TransmissionContainerPost{
					Checksum: 0x00,
				},
			},
			ExpectedChecksum: 0x1E,
		},
		{
			Name: "GetNamePayload checksum",
			Packet: TransmissionContainer{
				TransmissionContainerPre: TransmissionContainerPre{
					SrcMac: [6]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
					DstMac: [6]byte{0x54, 0x10, 0xEC, 0x03, 0x61, 0x50},
				},
				Packet: Packet{
					PacketPre: PacketPre{
						TAG:       0,
						Token:     uint32(0x00000000),
						CommandID: COMMANDID_GET_NAME,
					},
					payload: payload.EmptyPayload(),
					PacketPost: PacketPost{
						Checksum: 0x00,
					},
				},
				TransmissionContainerPost: TransmissionContainerPost{
					Checksum: 0x00,
				},
			},
			ExpectedChecksum: 0x4A,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(test *testing.T) {
			crc, err := testCase.Packet.getChecksum()
			if err != nil {
				test.Logf("unexpected exceptoin. %v", err)
				test.Fail()
			}
			if crc != testCase.ExpectedChecksum {
				test.Logf("Expected value: 0x%X, Actual value: 0x%X", testCase.ExpectedChecksum, crc)
				test.Fail()
			}
		})
	}
}
