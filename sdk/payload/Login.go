package payload

import (
	"bytes"
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

	usernameLength := byte(len(username))
	_, err := b.WriteString(string(usernameLength))
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

	payload := Payload{
		data:       buffBytes,
		dataLength: byte(len(buffBytes)),
	}

	login := &Login{
		Payload:  payload,
		username: username,
		password: password,
	}

	return login
}
