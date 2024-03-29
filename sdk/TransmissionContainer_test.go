package sdk

import (
	"bisecur/sdk/payload"
	"bisecur/sdk/payload/hcp"
	"bytes"
	"testing"
)

func TestTransmissionContainerEncode(t *testing.T) {
	testCases := []struct {
		Name                  string
		Request               TransmissionContainer
		ExpectedServerRequest string
	}{
		{
			Name: "Ping 1 request",
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
			Name: "Ping 2 request",
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
			Name: "Ping 3 request",
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
			Name: "Login request",
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
			Name: "Get Name request",
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
			Name: "Hm Get Transition request",
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
			Name: "Jcmp Get Values request",
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
			Name: "Jcmp Get Groups request",
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
		{
			Name: "Jcmp Get Users request",
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
					payload:    payload.JcmpPayload("{\"CMD\":\"GET_USERS\"}"),
					PacketPost: PacketPost{
						//Checksum: 0xA7,
					},
				},
				TransmissionContainerPost: TransmissionContainerPost{
					//Checksum: 0x33,
				},
			},
			ExpectedServerRequest: "0000000000065410EC8528BB001C0596833386067B22434D44223A224745545F5553455253227D58AE",
		},
		{ // 000000000006 5410EC036150 00 0B00 302B7D75 3300FF 8A 82
			Name: "Set State request",
			Request: TransmissionContainer{
				TransmissionContainerPre: TransmissionContainerPre{
					SrcMac: [6]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x06},
					DstMac: [6]byte{0x54, 0x10, 0xEC, 0x03, 0x61, 0x50},
				},
				Packet: Packet{
					PacketPre: PacketPre{
						TAG:       0x00,
						Token:     uint32(0x302B7D75),
						CommandID: COMMANDID_SET_STATE,
					},
					payload: payload.SetStatePayload(0x00),
					PacketPost: PacketPost{
						Checksum: 0x8A,
					},
				},
				TransmissionContainerPost: TransmissionContainerPost{
					Checksum: 0x82,
				},
			},
			ExpectedServerRequest: "0000000000065410EC036150000B00302B7D753300FF8A82",
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

			str := string(raw)
			expected := testCase.ExpectedServerRequest
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
			Name:         "Ping 1 Request",
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
			Name:         "Ping 2 Request",
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
			Name:         "Ping 3 Request",
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
			//  From my Own Device
			Name:         "Ping 4 Response",
			EncodedInput: "5410EC8528BB00000000000600090100000000808A7E",
			ExpectedDecodedInput: TransmissionContainer{
				TransmissionContainerPre: TransmissionContainerPre{
					SrcMac:     [6]byte{0x54, 0x10, 0xEC, 0x85, 0x28, 0xBB},
					DstMac:     [6]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x06},
					BodyLength: 9,
				},
				Packet: Packet{
					PacketPre: PacketPre{
						TAG:       1,
						Token:     uint32(0x00000000),
						CommandID: COMMANDID_PING_RESPONSE,
					},
					payload: payload.EmptyPayload(),
					PacketPost: PacketPost{
						Checksum: 0x8A,
					},
				},
				TransmissionContainerPost: TransmissionContainerPost{
					Checksum: 0x7E,
				},
			},
		},
		{
			Name:         "Login Request",
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
			Name:         "Login Failed Response",
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
		{
			Name:         "Permission Denied Response",
			EncodedInput: "5410EC036150000000000006000A007F162664010C36EB",
			ExpectedDecodedInput: TransmissionContainer{
				TransmissionContainerPre: TransmissionContainerPre{
					SrcMac:     [6]byte{0x54, 0x10, 0xEC, 0x03, 0x61, 0x50},
					DstMac:     [6]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x06},
					BodyLength: 0x0A,
				},
				Packet: Packet{
					PacketPre: PacketPre{
						TAG:       0x00,
						Token:     uint32(0x7F162664),
						CommandID: COMMANDID_ERROR,
					},
					payload: payload.ErrorPayload(payload.ERROR_PERMISSION_DENIED),
					PacketPost: PacketPost{
						Checksum: 0x36,
					},
				},
				TransmissionContainerPost: TransmissionContainerPost{
					Checksum: 0xEB,
				},
			},
		},
		{
			Name:         "Get Values Response",
			EncodedInput: "5410EC8528BB00000000000601CB0496833386867B223030223A312C223031223A302C223032223A302C223033223A302C223034223A302C223035223A302C223036223A302C223037223A302C223038223A302C223039223A302C223130223A302C223131223A302C223132223A302C223133223A302C223134223A302C223135223A302C223136223A302C223137223A302C223138223A302C223139223A302C223230223A302C223231223A302C223232223A302C223233223A302C223234223A302C223235223A302C223236223A302C223237223A302C223238223A302C223239223A302C223330223A302C223331223A302C223332223A302C223333223A302C223334223A302C223335223A302C223336223A302C223337223A302C223338223A302C223339223A302C223430223A302C223431223A302C223432223A302C223433223A302C223434223A302C223435223A302C223436223A302C223437223A302C223438223A302C223439223A302C223530223A302C223531223A302C223532223A302C223533223A302C223534223A302C223535223A302C223536223A302C223537223A302C223538223A302C223539223A302C223630223A302C223631223A302C223632223A302C223633223A38377D75F9",
			ExpectedDecodedInput: TransmissionContainer{
				TransmissionContainerPre: TransmissionContainerPre{
					SrcMac:     [6]byte{0x54, 0x10, 0xEC, 0x85, 0x28, 0xBB},
					DstMac:     [6]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x06},
					BodyLength: 0x01CB,
				},
				Packet: Packet{
					PacketPre: PacketPre{
						TAG:       0x04,
						Token:     uint32(0x96833386),
						CommandID: COMMANDID_JMCP_RESPONSE,
					},
					payload: payload.JcmpPayload("{\"00\":1,\"01\":0,\"02\":0,\"03\":0,\"04\":0,\"05\":0,\"06\":0,\"07\":0,\"08\":0,\"09\":0,\"10\":0,\"11\":0,\"12\":0,\"13\":0,\"14\":0,\"15\":0,\"16\":0,\"17\":0,\"18\":0,\"19\":0,\"20\":0,\"21\":0,\"22\":0,\"23\":0,\"24\":0,\"25\":0,\"26\":0,\"27\":0,\"28\":0,\"29\":0,\"30\":0,\"31\":0,\"32\":0,\"33\":0,\"34\":0,\"35\":0,\"36\":0,\"37\":0,\"38\":0,\"39\":0,\"40\":0,\"41\":0,\"42\":0,\"43\":0,\"44\":0,\"45\":0,\"46\":0,\"47\":0,\"48\":0,\"49\":0,\"50\":0,\"51\":0,\"52\":0,\"53\":0,\"54\":0,\"55\":0,\"56\":0,\"57\":0,\"58\":0,\"59\":0,\"60\":0,\"61\":0,\"62\":0,\"63\":87}"),
					PacketPost: PacketPost{
						Checksum: 0x75,
					},
				},
				TransmissionContainerPost: TransmissionContainerPost{
					Checksum: 0xF9,
				},
			},
		},
		{
			Name:         "Get Groups Response",
			EncodedInput: "5410EC8528BB000000000006003F0596833386865B7B226964223A302C226E616D65223A22676172617A73222C22706F727473223A5B7B226964223A302C2274797065223A317D5D7D5D88CF",
			ExpectedDecodedInput: TransmissionContainer{
				TransmissionContainerPre: TransmissionContainerPre{
					SrcMac:     [6]byte{0x54, 0x10, 0xEC, 0x85, 0x28, 0xBB},
					DstMac:     [6]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x06},
					BodyLength: 0x003F,
				},
				Packet: Packet{
					PacketPre: PacketPre{
						TAG:       0x05,
						Token:     uint32(0x96833386),
						CommandID: COMMANDID_JMCP_RESPONSE,
					},
					payload: payload.JcmpPayload("[{\"id\":0,\"name\":\"garazs\",\"ports\":[{\"id\":0,\"type\":1}]}]"),
					PacketPost: PacketPost{
						Checksum: 0x88,
					},
				},
				TransmissionContainerPost: TransmissionContainerPost{
					Checksum: 0xCF,
				},
			},
		},
		{
			Name:         "Login Response",
			EncodedInput: "5410EC8528BB000000000006000E0300000000900196833386748E",
			ExpectedDecodedInput: TransmissionContainer{
				TransmissionContainerPre: TransmissionContainerPre{
					SrcMac:     [6]byte{0x54, 0x10, 0xEC, 0x85, 0x28, 0xBB},
					DstMac:     [6]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x06},
					BodyLength: 0x000E,
				},
				Packet: Packet{
					PacketPre: PacketPre{
						TAG:       0x03,
						Token:     uint32(0x00000000),
						CommandID: COMMANDID_LOGIN_RESPONSE,
					},
					payload: payload.LoginResponsePayload(0x01, 0x96833386),
					PacketPost: PacketPost{
						Checksum: 0x74,
					},
				},
				TransmissionContainerPost: TransmissionContainerPost{
					Checksum: 0x8E,
				},
			},
		},
		{ // 5410EC036150 000000000006 0019 0000000000F0B800001801086C020000000000000000 50 93
			Name:         "Hm Get Transition Response",
			EncodedInput: "5410EC03615000000000000600190000000000F0B800001801086C0200000000000000005093",
			ExpectedDecodedInput: TransmissionContainer{
				TransmissionContainerPre: TransmissionContainerPre{
					SrcMac:     [6]byte{0x54, 0x10, 0xEC, 0x03, 0x61, 0x50},
					DstMac:     [6]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x06},
					BodyLength: 0x0019,
				},
				Packet: Packet{
					PacketPre: PacketPre{
						TAG:       0x00,
						Token:     uint32(0x00000000),
						CommandID: COMMANDID_HM_GET_TRANSITION_RESPONSE,
					},
					payload: &payload.HmGetTransitionResponse{
						Payload:               payload.MockPayload("B800001801086C020000000000000000"), // TODO not the most nice solution to validate but never want to generate this kind of package anyway
						StateInPercent:        92,
						DesiredStateInPercent: 0,
						Error:                 false,
						AutoClose:             false,
						DriveTime:             24,
						Gk:                    264,
						Hcp: &hcp.Hcp{
							PositionOpen:     false,
							PositionClose:    false,
							OptionRelais:     true,
							LightBarrier:     true,
							Error:            false,
							DrivingToClose:   true,
							Driving:          true,
							HalfOpened:       false,
							ForecastLeadTime: false,
							Learned:          true,
							NotReferenced:    false,
						},
						Exst: []byte{0, 0, 0, 0, 0, 0, 0, 0},
						//Time: time.UnixMilli(1701456598985267), // 2023-12-01T19:49:58.985267
						//IgnoreRetries: true,
					},
					PacketPost: PacketPost{
						Checksum: 0x50,
					},
				},
				TransmissionContainerPost: TransmissionContainerPost{
					Checksum: 0x93,
				},
			},
		},
		{
			// 5410EC8528BB 000000000006 000F 01 00000000 825410EC8528BB 4A 36
			Name:         "Get Mac Response",
			EncodedInput: "5410EC8528BB000000000006000F0100000000825410EC8528BB4A36",
			ExpectedDecodedInput: TransmissionContainer{
				TransmissionContainerPre: TransmissionContainerPre{
					SrcMac:     [6]byte{0x54, 0x10, 0xEC, 0x85, 0x28, 0xBB},
					DstMac:     [6]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x06},
					BodyLength: 0x0F,
				},
				Packet: Packet{
					PacketPre: PacketPre{
						TAG:       0x01,
						Token:     uint32(0x00000000),
						CommandID: COMMANDID_GET_MAC_RESPONSE,
					},
					payload: payload.GetMacResponsePayload([6]byte{0x54, 0x10, 0xEC, 0x85, 0x28, 0xBB}),
					PacketPost: PacketPost{
						Checksum: 0x4A,
					},
				},
				TransmissionContainerPost: TransmissionContainerPost{
					Checksum: 0x36,
				},
			},
		},
		{ // 5410EC8528BB 000000000006 0018 01 00000000 A6 426953656375722047617465776179 5E C3
			Name:         "Get Name Response",
			EncodedInput: "5410EC8528BB00000000000600180100000000A64269536563757220476174657761795EC3",
			ExpectedDecodedInput: TransmissionContainer{
				TransmissionContainerPre: TransmissionContainerPre{
					SrcMac:     [6]byte{0x54, 0x10, 0xEC, 0x85, 0x28, 0xBB},
					DstMac:     [6]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x06},
					BodyLength: 0x18,
				},
				Packet: Packet{
					PacketPre: PacketPre{
						TAG:       0x01,
						Token:     uint32(0x00000000),
						CommandID: COMMANDID_GET_NAME_RESPONSE,
					},
					payload: payload.GetNameResponsePayload("BiSecur Gateway"),
					PacketPost: PacketPost{
						Checksum: 0x5E,
					},
				},
				TransmissionContainerPost: TransmissionContainerPost{
					Checksum: 0xC3,
				},
			},
		},
		{ // 000000000005 5410EC8528BB 0021 00 884209B7 22 0A7465737A7466686E657631323334353637383930412161 F8 AA
			Name:         "Add new user",
			EncodedInput: "0000000000055410EC8528BB002100884209B7220A7465737A7466686E657631323334353637383930412161F8AA",
			ExpectedDecodedInput: TransmissionContainer{
				TransmissionContainerPre: TransmissionContainerPre{
					SrcMac:     [6]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x05},
					DstMac:     [6]byte{0x54, 0x10, 0xEC, 0x85, 0x28, 0xBB},
					BodyLength: 0x21,
				},
				Packet: Packet{
					PacketPre: PacketPre{
						TAG:       0x00,
						Token:     uint32(0x884209B7),
						CommandID: COMMANDID_ADD_USER,
					},
					payload: payload.LoginPayload("tesztfhnev", "1234567890A!a"),
					PacketPost: PacketPost{
						Checksum: 0xF8,
					},
				},
				TransmissionContainerPost: TransmissionContainerPost{
					Checksum: 0xAA,
				},
			},
		},
		{ // 5410EC8528BB 000000000005 000A 00 884209B7 A2 02 38 1B
			Name:         "Add new user response",
			EncodedInput: "5410EC8528BB000000000005000A00884209B7A202381B",
			ExpectedDecodedInput: TransmissionContainer{
				TransmissionContainerPre: TransmissionContainerPre{
					SrcMac:     [6]byte{0x54, 0x10, 0xEC, 0x85, 0x28, 0xBB},
					DstMac:     [6]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x05},
					BodyLength: 0x000A,
				},
				Packet: Packet{
					PacketPre: PacketPre{
						TAG:       0x00,
						Token:     uint32(0x884209B7),
						CommandID: COMMANDID_ADD_USER_RESPONSE,
					},
					payload: payload.AddUserResponsePayload(0x02),
					PacketPost: PacketPost{
						Checksum: 0x38,
					},
				},
				TransmissionContainerPost: TransmissionContainerPost{
					Checksum: 0x1B,
				},
			},
		},
		{ // 000000000005 5410EC8528BB 000A 00 884209B7 28 02 BE 2E
			Name:         "Get user rights",
			EncodedInput: "0000000000055410EC8528BB000A00884209B72802BE2E",
			ExpectedDecodedInput: TransmissionContainer{
				TransmissionContainerPre: TransmissionContainerPre{
					SrcMac:     [6]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x05},
					DstMac:     [6]byte{0x54, 0x10, 0xEC, 0x85, 0x28, 0xBB},
					BodyLength: 0x000A,
				},
				Packet: Packet{
					PacketPre: PacketPre{
						TAG:       0x00,
						Token:     uint32(0x884209B7),
						CommandID: COMMANDID_GET_USER_RIGHTS,
					},
					payload: payload.AddUserResponsePayload(0x02),
					PacketPost: PacketPost{
						Checksum: 0xBE,
					},
				},
				TransmissionContainerPost: TransmissionContainerPost{
					Checksum: 0x2E,
				},
			},
		},
		{ // 5410EC8528BB 000000000005 000A 00 884209B7 A8 02 3E 2E
			Name:         "Get user rights response",
			EncodedInput: "5410EC8528BB000000000005000A00884209B7A8023E2E",
			ExpectedDecodedInput: TransmissionContainer{
				TransmissionContainerPre: TransmissionContainerPre{
					SrcMac:     [6]byte{0x54, 0x10, 0xEC, 0x85, 0x28, 0xBB},
					DstMac:     [6]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x05},
					BodyLength: 0x000A,
				},
				Packet: Packet{
					PacketPre: PacketPre{
						TAG:       0x00,
						Token:     uint32(0x884209B7),
						CommandID: COMMANDID_GET_USER_RIGHTS_RESPONSE, // A8
					},
					payload: payload.GetUserRightsResponsePayload(0x02),
					PacketPost: PacketPost{
						Checksum: 0x3E,
					},
				},
				TransmissionContainerPost: TransmissionContainerPost{
					Checksum: 0x2E,
				},
			},
		},
		{ // 000000000005 5410EC8528BB 000A 00 884209B7 25 02 BB 28
			Name:         "Set user rights",
			EncodedInput: "0000000000055410EC8528BB000A00884209B72502BB28",
			ExpectedDecodedInput: TransmissionContainer{
				TransmissionContainerPre: TransmissionContainerPre{
					SrcMac:     [6]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x05},
					DstMac:     [6]byte{0x54, 0x10, 0xEC, 0x85, 0x28, 0xBB},
					BodyLength: 0x000A,
				},
				Packet: Packet{
					PacketPre: PacketPre{
						TAG:       0x00,
						Token:     uint32(0x884209B7),
						CommandID: COMMANDID_SET_USER_RIGHTS,
					},
					payload: payload.SetUserRightsPayload(0x02), // TODO what is this good for with only a user ID argument???
					PacketPost: PacketPost{
						Checksum: 0xBB,
					},
				},
				TransmissionContainerPost: TransmissionContainerPost{
					Checksum: 0x28,
				},
			},
		},
		{ // 5410EC8528BB 000000000005 000B 00 884209B7 A5 02 00 3C 8A
			Name:         "Set user rights response",
			EncodedInput: "5410EC8528BB000000000005000B00884209B7A502003C8A",
			ExpectedDecodedInput: TransmissionContainer{
				TransmissionContainerPre: TransmissionContainerPre{
					SrcMac:     [6]byte{0x54, 0x10, 0xEC, 0x85, 0x28, 0xBB},
					DstMac:     [6]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x05},
					BodyLength: 0x000B,
				},
				Packet: Packet{
					PacketPre: PacketPre{
						TAG:       0x00,
						Token:     uint32(0x884209B7),
						CommandID: COMMANDID_SET_USER_RIGHTS_RESPONSE,
					},
					payload: payload.SetUserRightsResponsePayload(0x02, 0x00), // TODO find out what is the second argument
					PacketPost: PacketPost{
						Checksum: 0x3C,
					},
				},
				TransmissionContainerPost: TransmissionContainerPost{
					Checksum: 0x8A,
				},
			},
		},
		{ // 000000000005 5410EC8528BB 000A 00 75730722 24 02 41 F1
			Name:         "Delete user",
			EncodedInput: "0000000000055410EC8528BB000A0075730722240241F1",
			ExpectedDecodedInput: TransmissionContainer{
				TransmissionContainerPre: TransmissionContainerPre{
					SrcMac:     [6]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x05},
					DstMac:     [6]byte{0x54, 0x10, 0xEC, 0x85, 0x28, 0xBB},
					BodyLength: 0x000A,
				},
				Packet: Packet{
					PacketPre: PacketPre{
						TAG:       0x00,
						Token:     uint32(0x75730722),
						CommandID: COMMANDID_REMOVE_USER,
					},
					payload: payload.RemoveUserPayload(0x02),
					PacketPost: PacketPost{
						Checksum: 0x41,
					},
				},
				TransmissionContainerPost: TransmissionContainerPost{
					Checksum: 0xF1,
				},
			},
		},
		{ // 5410EC8528BB 000000000005 000A 00 75730722 A4 02 C1 0F
			Name:         "Delete user response",
			EncodedInput: "5410EC8528BB000000000005000A0075730722A402C10F",
			ExpectedDecodedInput: TransmissionContainer{
				TransmissionContainerPre: TransmissionContainerPre{
					SrcMac:     [6]byte{0x54, 0x10, 0xEC, 0x85, 0x28, 0xBB},
					DstMac:     [6]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x05},
					BodyLength: 0x000A,
				},
				Packet: Packet{
					PacketPre: PacketPre{
						TAG:       0x00,
						Token:     uint32(0x75730722),
						CommandID: COMMANDID_REMOVE_USER_RESPONSE,
					},
					payload: payload.RemoveUserResponsePayload(0x02),
					PacketPost: PacketPost{
						Checksum: 0xC1,
					},
				},
				TransmissionContainerPost: TransmissionContainerPost{
					Checksum: 0x0F,
				},
			},
		},
		{ // 000000000005 5410EC8528BB 0017 00 75730722 45 0230393837363534333231765621 69 3B
			Name:         "Change user password",
			EncodedInput: "0000000000055410EC8528BB00170075730722450230393837363534333231765621693B",
			ExpectedDecodedInput: TransmissionContainer{
				TransmissionContainerPre: TransmissionContainerPre{
					SrcMac:     [6]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x05},
					DstMac:     [6]byte{0x54, 0x10, 0xEC, 0x85, 0x28, 0xBB},
					BodyLength: 0x0017,
				},
				Packet: Packet{
					PacketPre: PacketPre{
						TAG:       0x00,
						Token:     uint32(0x75730722),
						CommandID: COMMANDID_CHANGE_PASSWD,
					},
					payload: payload.ChangeUserPasswordPayload(0x02, "0987654321vV!"),
					PacketPost: PacketPost{
						Checksum: 0x69,
					},
				},
				TransmissionContainerPost: TransmissionContainerPost{
					Checksum: 0x3B,
				},
			},
		},
		{ // 5410EC8528BB 000000000005 0009 00 75730722 C5 DF BE
			Name:         "Change user password response",
			EncodedInput: "5410EC8528BB00000000000500090075730722C5DFBE",
			ExpectedDecodedInput: TransmissionContainer{
				TransmissionContainerPre: TransmissionContainerPre{
					SrcMac:     [6]byte{0x54, 0x10, 0xEC, 0x85, 0x28, 0xBB},
					DstMac:     [6]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x05},
					BodyLength: 0x0009,
				},
				Packet: Packet{
					PacketPre: PacketPre{
						TAG:       0x00,
						Token:     uint32(0x75730722),
						CommandID: COMMANDID_CHANGE_PASSWD_RESPONSE,
					},
					payload: payload.EmptyPayload(),
					PacketPost: PacketPost{
						Checksum: 0xDF,
					},
				},
				TransmissionContainerPost: TransmissionContainerPost{
					Checksum: 0xBE,
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(test *testing.T) {
			encodedInputBytes := []byte(testCase.EncodedInput)
			encodedInputBuffer := new(bytes.Buffer)
			_, err := encodedInputBuffer.Write(encodedInputBytes)
			if err != nil {
				test.Logf("Failed to write test data into buffer. %v", err)
				test.Fail()
				return
			}

			decoded, err := DecodeTransmissionContainer(encodedInputBuffer)
			if err != nil {
				test.Logf("Failed to decode request. %v", err)
				test.Fail()
				return
			}

			expected := &testCase.ExpectedDecodedInput
			if decoded == nil || !expected.Equal(decoded) {
				test.Logf("Expected value: %s\nActual value: %s", expected, decoded)
				test.Fail()
			} else {
				test.Logf("Decoded transmission container: %v", decoded)
			}
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
			Name: "Login checksum",
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
			Name: "Get Name Payload checksum",
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
