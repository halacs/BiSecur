package sdk

import (
	"reflect"
	"testing"
)

func TestGroupDecode(t *testing.T) {
	testCases := []struct {
		Name           string
		JsonInput      string
		ExpectedObject Groups
	}{
		{
			Name:      "One group only",
			JsonInput: "[{\"id\":0,\"name\":\"garazs\",\"ports\":[{\"id\":0,\"type\":1}]}]",
			ExpectedObject: Groups{
				{
					ID:   0,
					Name: "garazs",
					Ports: Ports{
						{
							ID:   0,
							Type: PORT_TYPE_IMPULS,
						},
					},
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(test *testing.T) {
			decoded, err := DecodeGroups(testCase.JsonInput)
			if err != nil {
				test.Logf("Unexcepted error happend. %v", err)
				test.Fail()
			}

			expected := testCase.ExpectedObject
			if !reflect.DeepEqual(expected, decoded) {
				test.Logf("Expected value: %+v, Actual value: %+v", expected, decoded)
				test.Fail()
			}
		})
	}
}
