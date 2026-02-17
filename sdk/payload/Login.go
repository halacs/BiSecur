package payload

import (
	"bytes"
	"halsecur/cli/utils"
)

type Login struct {
	Payload
	username string
	password string
}

func DecodeLoginPayload(payloadBytes []byte) (PayloadInterface, error) {
	usernameLength := payloadBytes[0]

	firstPasswordCharIndex := usernameLength + 1
	username := string(payloadBytes[1 : usernameLength+1])
	password := string(payloadBytes[firstPasswordCharIndex:])

	return LoginPayload(username, password), nil
}

func (l *Login) Encode() []byte {
	data := getHormanEncodedUsernamePassword(l.username, l.password)
	return data
}

func LoginPayload(username, password string) PayloadInterface {
	b := new(bytes.Buffer)

	userNameLength, err := utils.SafeLen(len(username))
	if err != nil {
		return nil
	}

	usernameLength := userNameLength
	_, err = b.WriteString(string(usernameLength))
	if err != nil {
		panic("???")
	}

	_, err = b.WriteString(username)
	if err != nil {
		panic("???")
	}

	_, err = b.WriteString(password)
	if err != nil {
		panic("???")
	}

	buffBytes := b.Bytes()

	dataLength, err := utils.SafeLen(len(buffBytes))
	if err != nil {
		return nil
	}

	payload := Payload{
		data:       buffBytes,
		dataLength: dataLength,
	}

	login := &Login{
		Payload:  payload,
		username: username,
		password: password,
	}

	return login
}
