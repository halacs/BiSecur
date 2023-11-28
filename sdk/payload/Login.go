package payload

import (
	"bytes"
)

type Login struct {
	Payload
	username string
	password string
}

func (l *Login) Encode() []byte {
	return getHormanEncodedUsernamePassword(l.username, l.password)
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
		data: buffBytes,
	}

	login := &Login{
		Payload:  payload,
		username: username,
		password: password,
	}

	return login
}
