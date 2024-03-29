package sdk

import (
	"encoding/json"
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
	json, err := json.Marshal(u)
	if err != nil {
		panic(err)
	}
	return string(json)
}

func (u *Users) String() string {
	json, err := json.Marshal(u)
	if err != nil {
		panic(err)
	}
	return string(json)
}
