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
			JsonInput: "[{\"id\":0,\"name\":\"admin\",\"isAdmin\":true,\"Groups\":[]},{\"id\":1,\"name\":\"app\",\"isAdmin\":false,\"Groups\":[0]}]",
			ExpectedObject: Users{
				User{
					ID:      0,
					Name:    "admin",
					IsAdmin: true,
					Groups:  []int{},
				},
				User{
					ID:      1,
					Name:    "app",
					IsAdmin: false,
					Groups:  []int{0},
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
