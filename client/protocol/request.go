package protocol

type Request struct {
	Agencia int
}

const (
	RequestIdentifier = 2
)

func (request *Request) ToBytes() []byte {
	bytes := make([]byte, 0, 2)
	// Agregamos al principio un u8 que le diga al servidor que este mensaje es un pedido de ganadores
	bytes = append(bytes, byte(RequestIdentifier))
	// Agrego el número de agencia para que me sepan identificar
	bytes = append(bytes, byte(request.Agencia))
	return bytes
}
