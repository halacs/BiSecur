package sdk

import (
	"encoding/json"
	"fmt"
)

type Users []User

// Example:
// [{"id":0,"name":"admin","isAdmin":true,"Groups":[]},{"id":1,"name":"app","isAdmin":false,"Groups":[0]}]
type User struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	IsAdmin bool   `json:"isAdmin"`
	Groups  []int  `json:"Groups"`
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
	return fmt.Sprintf("[ID=%d, Name=\"%s\", IsAdmin=%t, Groups:%v]", u.ID, u.Name, u.IsAdmin, u.Groups)
}

func (users Users) String() string {
	s := ""
	for _, u := range users {
		s = s + u.String()
	}
	return s
}
