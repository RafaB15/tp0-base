package common

import (
	"encoding/csv"
	"errors"
	"io"
	"net"
	"os"

	p "github.com/7574-sistemas-distribuidos/docker-compose-init/client/protocol"
	"github.com/op/go-logging"
)

const (
	Tries = 2
)

var log = logging.MustGetLogger("log")

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            string
	ServerAddress string
	BatchSize     int
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

// Send a batch of bets through the socket. I serializes them, sends them and then waits for confirmation.
func (c *Client) SendBetBatch(betBatch []*p.Bet) error {
	serializedBets := p.SerializeBetBatch(betBatch)
	err := p.WriteExact(c.conn, serializedBets)
	if err != nil {
		return err
	}

	confirmation, err := p.ReadExact(c.conn, 1)
	if err != nil {
		return err
	}

	if int(confirmation[0]) != len(betBatch) {
		return errors.New("algunas de las consultas no pudieron ser procesadas")
	}

	return nil
}

// Attempts to send a batch of messages an specified amount of times.
func (c *Client) SendBetBatchWithTries(tries int, betBatch []*p.Bet) {
	currentTries := 0
	for currentTries < tries {
		log.Infof("action: batch_enviado | result: in_progress | client_id: %v | cantidad_a_enviar: %d", c.config.ID, len(betBatch))
		err := c.SendBetBatch(betBatch)
		if err == nil {
			log.Infof("action: batch_enviado | result: success | client_id: %v | cantidad_enviada: %d", c.config.ID, len(betBatch))
			return
		}
		log.Errorf("action: batch_enviado | result: fail | client_id: %v | error: %v", c.config.ID, err)
		tries += 1
	}
}

// SendBet takes a bet and sends it to the server. I wait to receive a confimfmation.
func (c *Client) SendBets(pathBets string, agencia int) {
	file, err := os.Open(pathBets)
	if err != nil {
		log.Errorf("action: open_file | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
		return
	}

	defer file.Close()

	reader := csv.NewReader(file)
	bets := []*p.Bet{}

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Errorf("action: read_bet | result: fail | client_id: %v | error: %v",
				c.config.ID,
				err,
			)
			return
		}

		if len(line) != 5 {
			log.Errorf("action: read_bet | result: fail | client_id: %v | error: Insufficient data on line",
				c.config.ID,
			)
			continue
		}

		bet, err := p.NewBet(
			agencia,
			line[0],
			line[1],
			line[2],
			line[3],
			line[4],
		)

		if err != nil {
			log.Errorf("action: read_bet | result: fail | client_id: %v | error: %v",
				c.config.ID,
				err,
			)
			continue
		}

		bets = append(bets, bet)

		if len(bets) >= c.config.BatchSize {
			c.SendBetBatchWithTries(Tries, bets)
			bets = []*p.Bet{}
		}
	}

	if len(bets) > 0 {
		c.SendBetBatchWithTries(Tries, bets)
	}
}

func (c *Client) GetWinners(agencia int) {
	request := p.Request{
		Agencia: agencia,
	}

	serializedRequest := request.ToBytes()
	err := p.WriteExact(c.conn, serializedRequest)
	if err != nil {
		log.Errorf("action: serialize_request | result: fail | error %v", err)
		return
	}

	winners, err := p.DeserializeWinners(c.conn)
	if err != nil {
		log.Errorf("action: consulta_ganadores | result: fail | error %v", err)
		return
	}

	log.Infof("action: consulta_ganadores | result: success | cant_ganadores %d", len(winners.Documentos))
}

func (c *Client) GetBetsResults(pathBets string, agencia int) {
	err := c.createClientSocket()
	defer c.conn.Close()
	if err != nil {
		log.Errorf("action: connect_to_server | result: fail | error: %v", err)
	}

	c.SendBets(pathBets, agencia)

	c.GetWinners(agencia)
}

// Closes the connection of the client
func (c *Client) Shutdown() {
	if c.conn != nil {
		c.conn.Close()
		log.Infof("action: connection_closed | result: success | client_id: %v", c.config.ID)
	}
}
