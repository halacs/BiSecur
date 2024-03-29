package sdk

import "encoding/json"

type Values map[string]byte

func DecodeValues(jsonStr string) (Values, error) {
	var data Values

	err := json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
