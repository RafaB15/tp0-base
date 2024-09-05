package protocol

import (
	"encoding/binary"
	"errors"
	"strconv"
)

const (
	AgenciaLength        = 1
	NombreLength         = 1
	ApellidoLength       = 1
	DocumentoLength      = 8
	NacimientoLength     = 10
	NumeroLength         = 2
	BetMessageIdentifier = 1
)

type Bet struct {
	Agencia    int
	Nombre     string
	Apellido   string
	Documento  string
	Nacimiento string
	Numero     int
}

func (bet *Bet) ToBytes() []byte {
	// Calculo la longitud de los campos de longitud variables
	nombreLen := len(bet.Nombre)
	apellidoLen := len(bet.Apellido)

	// Creamos un slice de bytes en donde iremos poniendo los campos
	bytes := make([]byte, 0, AgenciaLength+NombreLength+nombreLen+ApellidoLength+apellidoLen+DocumentoLength+NacimientoLength+NumeroLength)

	// Agregamos un u8 que representa el número de agencia
	bytes = append(bytes, byte(bet.Agencia))

	// Agregamos la longitud del nombre como u8 y el nombre
	bytes = append(bytes, byte(nombreLen))
	bytes = append(bytes, []byte(bet.Nombre)...)

	// Agregamos la longitud del apellido como u8 y el apellido
	bytes = append(bytes, byte(apellidoLen))
	bytes = append(bytes, []byte(bet.Apellido)...)

	// Agregamos un string con el documento. Son 8 caracteres fijos según mi protocolo.
	bytes = append(bytes, []byte(bet.Documento)...)

	// Agregamos fecha en formato específico (fixed format 2000-04-15)
	bytes = append(bytes, []byte(bet.Nacimiento)...)

	// Agregamos el número como un u16
	numBytes := make([]byte, NumeroLength)
	binary.BigEndian.PutUint16(numBytes, uint16(bet.Numero))
	bytes = append(bytes, numBytes...)

	// Obtenemos la longitud total y la ponemos al principio de todo.
	totalLength := len(bytes)
	finalBytes := make([]byte, 2+totalLength)
	binary.BigEndian.PutUint16(finalBytes, uint16(totalLength))
	copy(finalBytes[2:], bytes)

	return finalBytes
}

func SerializeBetBatch(betBatch []*Bet) []byte {
	serialized := []byte{BetMessageIdentifier, byte(len(betBatch))}

	for _, bet := range betBatch {
		serialized = append(serialized, bet.ToBytes()...)
	}

	return serialized
}

func NewBet(agencia int, nombre string, apellido string, documento string, nacimiento string, numeroStr string) (*Bet, error) {
	if nombre == "" {
		return nil, errors.New("el nombre no puede estar vacío")
	}

	if apellido == "" {
		return nil, errors.New("el apellido no puede estar vacío")
	}

	if documento == "" || len(documento) != 8 {
		return nil, errors.New("el documento tiene que tener 8 caracteres")
	}

	if nacimiento == "" {
		return nil, errors.New("el nacimiento no puede estar vacío")
	}

	numero, err := strconv.Atoi(numeroStr)
	if err != nil {
		return nil, errors.New("número inválido")
	}

	bet := &Bet{
		Agencia:    agencia,
		Nombre:     nombre,
		Apellido:   apellido,
		Documento:  documento,
		Nacimiento: nacimiento,
		Numero:     numero,
	}

	return bet, nil

}
