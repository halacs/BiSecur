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
					//payload: payload.LoginPayload("thomas", "aaabbbccc"),
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

/*
func TestCommandRequestDecoding(t *testing.T) {
	testCases := []struct {
		Name                   string
		CommandRequest         string
		ExpectedCommandRequest CommandRequest
	}{
		{
			Name:           "CommandCodec12Request1",
			CommandRequest: "000000000000000F0C010500000007676574696E666F0100004312",
			ExpectedCommandRequest: CommandRequest{
				commandRequestPre: commandRequestPre{
					Preamble:         0,
					DataSize:         0x0F,
					CodecID:          0x0C,
					CommandQuantity1: 1,
					Type:             0x05,
					CommandSize:      0x00000007,
				},
				Command: []byte("getinfo"),
				commandRequestPost: commandRequestPost{
					CommandQuantity2: 0x01,
					Checksum:              0x4312,
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(test *testing.T) {
			rawCommandRequest, err := hex.DecodeString(testCase.CommandRequest)
			if err != nil {
				test.Logf("Failed to decode client string to byte array. %v", err)
				test.Fail()
			}

			decoded, err := DecodeCommandRequest(&rawCommandRequest)
			if err != nil {
				test.Logf("Failed to decode command request. %v", err)
				test.Fail()
			}

			if !reflect.DeepEqual(decoded, testCase.ExpectedCommandRequest) {
				test.Logf("Expected value: %v, Actual value: %v", testCase.ExpectedCommandRequest, decoded)
				test.Fail()
			}
		})
	}
}
*/

/*
func TestCommandResponseDecode(t *testing.T) {
	testCases := []struct {
		Name                    string
		ErrorCase               bool
		ExpectedErrorMessage    string
		ClientResponse          string
		ExpectedDecodedResponse CommandResponse
	}{
		{
			Name: "CommandCodec12GetInfoResponse",

			ClientResponse:       "00000000000000900C010600000088494E493A323031392F372F323220373A3232205254433A323031392F372F323220373A3533205253543A32204552523A312053523A302042523A302043463A302046473A3020464C3A302054553A302F302055543A3020534D533A30204E4F4750533A303A3330204750533A31205341543A302052533A332052463A36352053463A31204D443A30010000C78F",
			ErrorCase:            false,
			ExpectedErrorMessage: "",
			ExpectedDecodedResponse: CommandResponse{
				commandResponsePre: commandResponsePre{
					Preamble:          0x00000000,
					DataSize:          0x90,
					CodecID:           0x0C,
					ResponseQuantity1: 0x01,
					Type:              0x06,
					ResponseSize:      0x88,
				},
				Response: []byte("INI:2019/7/22 7:22 RTC:2019/7/22 7:53 RST:2 ERR:1 SR:0 BR:0 CF:0 FG:0 FL:0 TU:0/0 UT:0 SMS:0 NOGPS:0:30 GPS:1 SAT:0 RS:3 RF:65 SF:1 MD:0"),
				commandResponsePost: commandResponsePost{
					ResponseQuantity2: 0x01,
					Checksum:               0xC78F,
				},
			},
		},
		{
			Name:                 "CommandCodec12GetIoResponse",
			ErrorCase:            false,
			ExpectedErrorMessage: "",
			ClientResponse:       "00000000000000370C01060000002F4449313A31204449323A30204449333A302041494E313A302041494E323A313639323420444F313A3020444F323A3101000066E3",
			ExpectedDecodedResponse: CommandResponse{
				commandResponsePre: commandResponsePre{
					Preamble:          0x00000000,
					DataSize:          0x37,
					CodecID:           0x0C,
					ResponseQuantity1: 0x01,
					Type:              0x06,
					ResponseSize:      0x2F,
				},
				Response: []byte("DI1:1 DI2:0 DI3:0 AIN1:0 AIN2:16924 DO1:0 DO2:1"),
				commandResponsePost: commandResponsePost{
					ResponseQuantity2: 0x01,
					Checksum:               0x66E3,
				},
			},
		},
		{
			Name:                    "CommandCodec12GetIoResponseWrongPreamble",
			ErrorCase:               true,
			ExpectedErrorMessage:    "wrong preamble: 0x10000000",
			ClientResponse:          "10000000000000370C01060000002F4449313A31204449323A30204449333A302041494E313A302041494E323A313639323420444F313A3020444F323A3101000066E3",
			ExpectedDecodedResponse: CommandResponse{},
		},
		{
			Name:                    "CommandCodec12GetIoResponseWrongCrc",
			ErrorCase:               true,
			ExpectedErrorMessage:    "wrong Checksum! Calculated: 66e3 Received: 66e4",
			ClientResponse:          "00000000000000370C01060000002F4449313A31204449323A30204449333A302041494E313A302041494E323A313639323420444F313A3020444F323A3101000066E4",
			ExpectedDecodedResponse: CommandResponse{},
		},
		{
			Name:                    "CommandCodec12GetIoResponseWrongType",
			ErrorCase:               true,
			ExpectedErrorMessage:    "wrong type: 0x5",
			ClientResponse:          "00000000000000370C01050000002F4449313A31204449323A30204449333A302041494E313A302041494E323A313639323420444F313A3020444F323A3101000066E3",
			ExpectedDecodedResponse: CommandResponse{},
		},
		{
			Name:                    "CommandCodec12GetIoResponseTooShortMessage",
			ErrorCase:               true,
			ExpectedErrorMessage:    "wrong CodecID: 0x33",
			ClientResponse:          "000000000176000f3335303432343036333831373336338e01000001839ed29768000b5629e81c5451d0000000000000000000000b000500500000150400c800004502001d00000500422e9b0018000000cd13f000ce005d00430fd4000100f10000547e0000000001", // modified Codec8 payload package
			ExpectedDecodedResponse: CommandResponse{},
		},
		{
			Name:                    "Try to decode something different from a command response",
			ErrorCase:               true,
			ExpectedErrorMessage:    "only 7 bytes received. Probably not a teltonika command response packet",
			ClientResponse:          "0005cafe017601", // server acknowledgement package of Codec8 or Codec8 Extended
			ExpectedDecodedResponse: CommandResponse{},
		},
	}

	// Run all natsio cases as a separated network connection
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(test *testing.T) {
			rawClientResponse, err := hex.DecodeString(testCase.ClientResponse)
			if err != nil {
				test.Logf("Failed to decode client string to byte array. %v", err)
				test.Fail()
			}

			decoded, err := DecodeCommandResponse(&rawClientResponse)
			if testCase.ErrorCase {
				if err == nil {
					test.Logf("This is an error case but there is no error.")
					test.Fail()
				}
				if err.Error() != testCase.ExpectedErrorMessage {
					test.Logf("Expected error message: %v, Actual error message: %v", testCase.ExpectedErrorMessage, err.Error())
					test.Fail()
				}
				return
			}

			if err != nil {
				test.Logf("Failed to decode client request. %v", err)
				test.Fail()
			}

			if !reflect.DeepEqual(decoded, testCase.ExpectedDecodedResponse) {
				test.Logf("Expected value: %v, Actual value: %v", testCase.ExpectedDecodedResponse, decoded)
				test.Fail()
			}
		})
	}
}
*/

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
