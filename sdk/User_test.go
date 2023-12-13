package sdk

import (
	"reflect"
	"testing"
)

func TestUserDecode(t *testing.T) {
	testCases := []struct {
		Name           string
		JsonInput      string
		ExpectedObject Users
	}{
		{
			Name:      "Two users available",
			JsonInput: "[{\"id\":0,\"name\":\"admin\",\"isAdmin\":true,\"groups\":[]},{\"id\":1,\"name\":\"app\",\"isAdmin\":false,\"groups\":[0]}]",
			ExpectedObject: Users{
				User{
					ID:      0,
					Name:    "admin",
					IsAdmin: true,
					groups:  nil,
				},
				User{
					ID:      1,
					Name:    "app",
					IsAdmin: false,
					groups:  nil,
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(test *testing.T) {
			decoded, err := DecodeUsers(testCase.JsonInput)
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
