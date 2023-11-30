package payload

import "encoding/hex"

const (
	ERROR_COMMAND_NOT_FOUND   = 0
	ERROR_INVALID_PROTOCOL    = 1
	ERROR_LOGIN_FAILED        = 2
	ERROR_INVALID_TOKEN       = 3
	ERROR_USER_ALREADY_EXISTS = 4
	ERROR_NO_EMPTY_USER_SLOT  = 5
	ERROR_INVALID_PASSWORD    = 6
	ERROR_INVALID_USERNAME    = 7
	ERROR_USER_NOT_FOUND      = 8
	ERROR_PORT_NOT_FOUND      = 9
	ERROR_PORT_ERROR          = 10
	ERROR_GATEWAY_BUSY        = 11
	ERROR_PERMISSION_DENIED   = 12
	ERROR_NO_EMPTY_GROUP_SLOT = 13
	ERROR_GROUP_NOT_FOUND     = 14
	ERROR_INVALID_PAYLOAD     = 15
	ERROR_OUT_OF_RANGE        = 16
	ERROR_ADD_PORT_ERROR      = 17
	ERROR_NO_EMPTY_PORT_SLOT  = 18
	ERROR_ADAPTER_BUSY        = 19
)

type Error struct {
	Payload
}

func DecodeErrorPacket(payloadBytes []byte) (PayloadInterface, error) {
	return &Error{
		Payload{
			data: payloadBytes, // error code
		},
	}, nil
}

func (p *Error) Encode() []byte {
	encodedBytes := make([]byte, 2)
	hex.Encode(encodedBytes, p.ToByteArray())
	return encodedBytes
}

func ErrorPayload(errorCode byte) PayloadInterface {
	return &Error{
		Payload{
			data: []byte{errorCode},
		},
	}
}

func (e *Error) GetErrorCode() byte {
	return e.ToByteArray()[0]
}

func (e *Error) String() string {
	switch e.GetErrorCode() {
	case ERROR_COMMAND_NOT_FOUND:
		return "COMMAND_NOT_FOUND"
	case ERROR_INVALID_PROTOCOL:
		return "INVALID_PROTOCOL"
	case ERROR_LOGIN_FAILED:
		return "LOGIN_FAILED"
	case ERROR_INVALID_TOKEN:
		return "INVALID_TOKEN"
	case ERROR_USER_ALREADY_EXISTS:
		return "USER_ALREADY_EXISTS"
	case ERROR_NO_EMPTY_USER_SLOT:
		return "NO_EMPTY_USER_SLOT"
	case ERROR_INVALID_PASSWORD:
		return "INVALID_PASSWORD"
	case ERROR_INVALID_USERNAME:
		return "INVALID_USERNAME"
	case ERROR_USER_NOT_FOUND:
		return "USER_NOT_FOUND"
	case ERROR_PORT_NOT_FOUND:
		return "PORT_NOT_FOUND"
	case ERROR_PORT_ERROR:
		return "PORT_ERROR"
	case ERROR_GATEWAY_BUSY:
		return "GATEWAY_BUSY"
	case ERROR_PERMISSION_DENIED:
		return "PERMISSION_DENIED"
	case ERROR_NO_EMPTY_GROUP_SLOT:
		return "NO_EMPTY_GROUP_SLOT"
	case ERROR_GROUP_NOT_FOUND:
		return "GROUP_NOT_FOUND"
	case ERROR_INVALID_PAYLOAD:
		return "INVALID_PAYLOAD"
	case ERROR_OUT_OF_RANGE:
		return "OUT_OF_RANGE"
	case ERROR_ADD_PORT_ERROR:
		return "ADD_PORT_ERROR"
	case ERROR_NO_EMPTY_PORT_SLOT:
		return "NO_EMPTY_PORT_SLOT"
	case ERROR_ADAPTER_BUSY:
		return "ADAPTER_BUSY"
	}

	return "unknown"
}
