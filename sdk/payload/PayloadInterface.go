package payload

type PayloadInterface interface {
	ToByteArray() []byte
	Encode() []byte
	Length() int
}
