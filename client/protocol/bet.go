package protocol

import (
	"encoding/binary"
)

const (
	AgenciaLength    = 1
	NombreLength     = 1
	ApellidoLength   = 1
	DocumentoLength  = 8
	NacimientoLength = 10
	NumeroLength     = 2
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
