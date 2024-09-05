package protocol

import (
	"encoding/binary"
	"net"
)

type Winners struct {
	Documentos []string
}

func DeserializeWinners(socket net.Conn) (*Winners, error) {
	var winners Winners

	var numDocuments uint16
	err := binary.Read(socket, binary.BigEndian, &numDocuments)
	if err != nil {
		return nil, err
	}

	for i := 0; i < int(numDocuments); i++ {
		document, err := ReadExact(socket, 8)
		if err != nil {
			return nil, err
		}
		winners.Documentos = append(winners.Documentos, string(document))
	}

	return &winners, nil
}
