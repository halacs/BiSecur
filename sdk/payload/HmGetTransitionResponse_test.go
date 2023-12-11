package payload

import (
	"bisecure/sdk/payload/hcp"
	"encoding/hex"
	"testing"
	"time"
)

func TestTransmissionContainerDecode(t *testing.T) {
	testCases := []struct {
		Name                 string
		EncodedInput         string
		ExpectedDecodedInput HmGetTransitionResponse
	}{
		{
			Name:         "Open Door Transition decode",
			EncodedInput: "C8C80000010801020000000000000000",
			ExpectedDecodedInput: HmGetTransitionResponse{
				StateInPercent:        100,
				DesiredStateInPercent: 100,
				Error:                 false,
				AutoClose:             false,
				DriveTime:             0,
				Gk:                    264,
				Hcp: &hcp.Hcp{
					PositionOpen:     true,
					PositionClose:    false,
					OptionRelais:     false,
					LightBarrier:     false,
					Error:            false,
					DrivingToClose:   false,
					Driving:          false,
					HalfOpened:       false,
					ForecastLeadTime: false,
					Learned:          true,
					NotReferenced:    false,
				},
				Exst: []byte{0, 0, 0, 0, 0, 0, 0, 0},
				Time: time.UnixMilli(1701456598985267), // 2023-12-01T19:49:58.985267
				//IgnoreRetries: true,
			},
		},
		{
			Name:         "Closed Door Transition decode",
			EncodedInput: "00000000010802020000000000000000",
			ExpectedDecodedInput: HmGetTransitionResponse{
				StateInPercent:        0,
				DesiredStateInPercent: 0,
				Error:                 false,
				AutoClose:             false,
				DriveTime:             0,
				Gk:                    264,
				Hcp: &hcp.Hcp{
					PositionOpen:     false,
					PositionClose:    true,
					OptionRelais:     false,
					LightBarrier:     false,
					Error:            false,
					DrivingToClose:   false,
					Driving:          false,
					HalfOpened:       false,
					ForecastLeadTime: false,
					Learned:          true,
					NotReferenced:    false,
				},
				Exst: []byte{0, 0, 0, 0, 0, 0, 0, 0},
				//Time: time.UnixMilli(1701456598985267), // 2023-12-01T19:49:58.985267
				//IgnoreRetries: true,
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

			decoded, err := DecodeHmGetTransitionResponsePayload(encodedInputBytes)
			if err != nil {
				test.Logf("Failed to decode request. %v", err)
				test.Fail()
				return
			}

			expected := &testCase.ExpectedDecodedInput

			if !expected.Equal(decoded.(*HmGetTransitionResponse)) {
				test.Logf("Expected value: %v, Actual value: %v", expected, decoded)
				test.Fail()
			}
			test.Logf("Decoded transmission container: %v", decoded)
		})
	}
}
