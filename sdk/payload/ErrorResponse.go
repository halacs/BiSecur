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

type ErrorResponse struct {
	Payload
}

func DecodeErrorPayload(payloadBytes []byte) (PayloadInterface, error) {
	return &ErrorResponse{
		Payload{
			data: payloadBytes, // error code
		},
	}, nil
}

func (p *ErrorResponse) Encode() []byte {
	encodedBytes := make([]byte, 2)
	hex.Encode(encodedBytes, p.ToByteArray())
	return encodedBytes
}

func ErrorPayload(errorCode byte) PayloadInterface {
	return &ErrorResponse{
		Payload{
			data: []byte{errorCode},
		},
	}
}

func (e *ErrorResponse) GetErrorCode() byte {
	return e.ToByteArray()[0]
}

func (e *ErrorResponse) Equal(other *ErrorResponse) bool {
	return e.GetErrorCode() == other.GetErrorCode()
}

// Implement error interface
func (e *ErrorResponse) Error() string {
	return e.String()
}

func (e *ErrorResponse) String() string {
	switch e.GetErrorCode() {
	case ERROR_COMMAND_NOT_FOUND:
		return "COMMAND NOT FOUND"
	case ERROR_INVALID_PROTOCOL:
		return "INVALID PROTOCOL"
	case ERROR_LOGIN_FAILED:
		return "LOGIN FAILED"
	case ERROR_INVALID_TOKEN:
		return "INVALID TOKEN"
	case ERROR_USER_ALREADY_EXISTS:
		return "USER ALREADY EXISTS"
	case ERROR_NO_EMPTY_USER_SLOT:
		return "NO EMPTY USER SLOT"
	case ERROR_INVALID_PASSWORD:
		return "INVALID PASSWORD"
	case ERROR_INVALID_USERNAME:
		return "INVALID USERNAME"
	case ERROR_USER_NOT_FOUND:
		return "USER NOT FOUND"
	case ERROR_PORT_NOT_FOUND:
		return "PORT NOT FOUND"
	case ERROR_PORT_ERROR:
		return "PORT ERROR"
	case ERROR_GATEWAY_BUSY:
		return "GATEWAY BUSY"
	case ERROR_PERMISSION_DENIED:
		return "PERMISSION DENIED"
	case ERROR_NO_EMPTY_GROUP_SLOT:
		return "NO EMPTY GROUP SLOT"
	case ERROR_GROUP_NOT_FOUND:
		return "GROUP NOT FOUND"
	case ERROR_INVALID_PAYLOAD:
		return "INVALID PAYLOAD"
	case ERROR_OUT_OF_RANGE:
		return "OUT OF RANGE"
	case ERROR_ADD_PORT_ERROR:
		return "ADD PORT ERROR"
	case ERROR_NO_EMPTY_PORT_SLOT:
		return "NO EMPTY PORT SLOT"
	case ERROR_ADAPTER_BUSY:
		return "ADAPTER BUSY"
	}

	return "unknown"
}
