package sdk

import (
	"reflect"
	"testing"
)

func TestGroupDecode(t *testing.T) {
	testCases := []struct {
		Name           string
		JsonInput      string
		ExpectedObject []Group
	}{
		{
			Name:      "One group only",
			JsonInput: "[{\"id\":0,\"name\":\"garazs\",\"ports\":[{\"id\":0,\"type\":1}]}]",
			ExpectedObject: []Group{
				{
					ID:   0,
					Name: "garazs",
					Ports: []Port{
						{
							ID:   0,
							Type: 1,
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
