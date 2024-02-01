package sdk

import (
	"encoding/json"
)

type Groups []Group

// Example:
// [{"id":0,"name":"garazs","ports":[{"id":0,"type":1}]}]
type Group struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Ports Ports  `json:"ports"`
}

func DecodeGroups(jsonStr string) (Groups, error) {
	var data Groups

	err := json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (groups Groups) String() string {
	json, err := json.Marshal(groups)
	if err != nil {
		panic(err)
	}
	return string(json)
}
