package common

import (
	"net"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/protocol"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            string
	ServerAddress string
}

// Client Entity that encapsulates how
type Client struct {
	config ClientConfig
	conn   net.Conn
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig) *Client {
	client := &Client{
		config: config,
	}
	return client
}

// CreateClientSocket Initializes client socket. In case of
// failure, error is printed in stdout/stderr and exit 1
// is returned
func (c *Client) createClientSocket() error {
	conn, err := net.Dial("tcp", c.config.ServerAddress)
	if err != nil {
		log.Criticalf(
			"action: connect | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
	}
	c.conn = conn
	return nil
}

// SendBetConn sends a bet to the server. If the initial write does not send all the bytes,
// it continues writing the remaining bytes until all bytes are sent.
func (c *Client) SendBetConn(bet *protocol.Bet) bool {
	serializedBet := bet.ToBytes()
	err := writeExact(c.conn, serializedBet)
	if err != nil {
		log.Errorf("action: apuesta_enviada | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err)
	}

	// Espero confirmación del servidor
	confirmation, err := readExact(c.conn, 1)
	if err != nil {
		log.Errorf("action: read_confirmation | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
		return false
	}

	if confirmation[0] == 1 {
		log.Infof("action: apuesta_enviada | result: success | dni: %v | numero: %v", bet.Documento, bet.Numero)
		return true
	} else {
		log.Infof("action: apuesta_enviada | result: fail | dni: %v | numero: %v", bet.Documento, bet.Numero)
		return false
	}
}

// SendBet takes a bet and sends it to the server. I wait to receive a confimfmation.
func (c *Client) SendBet(bet *protocol.Bet) {
	// Create the connection the server. Send an
	c.createClientSocket()
	c.SendBetConn(bet)
	c.conn.Close()
}

func (c *Client) Shutdown() {
	if c.conn != nil {
		c.conn.Close()
		log.Infof("action: connection_closed | result: success | client_id: %v", c.config.ID)
	}
}
