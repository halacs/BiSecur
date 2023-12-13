package sdk

import (
	"encoding/json"
	"fmt"
)

type Users []User

// Example:
// [{"id":0,"name":"admin","isAdmin":true,"groups":[]},{"id":1,"name":"app","isAdmin":false,"groups":[0]}]
type User struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	IsAdmin bool   `json:"isAdmin"`
	groups  []int  `json:"groups"`
}

func DecodeUsers(jsonStr string) (Users, error) {
	var data Users

	err := json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (u *User) String() string {
	return fmt.Sprintf("[ID=%d, Name=\"%s\", IsAdmin=%t, Groups:%v]", u.ID, u.Name, u.IsAdmin, u.groups)
}

func (users Users) toString() string {
	s := ""
	for _, u := range users {
		s = s + u.String()
	}
	return s
}
