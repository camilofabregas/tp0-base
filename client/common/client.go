package common

import (
	"bufio"
	"fmt"
	"net"
	"time"
	"os"
    "os/signal"
    "syscall"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            string
	ServerAddress string
	LoopAmount    int
	LoopPeriod    time.Duration
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

// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) StartClientLoop() {
	sigc := make(chan os.Signal, 1)
    signal.Notify(sigc, syscall.SIGTERM)

	go func() {
        <-sigc
        log.Infof("action: sigterm_received | result: in_progress | client_id: %v | msg: Starting client shutdown", c.config.ID)
        if c.conn != nil {
            c.conn.Close()
        }

		log.Infof("action: shutdown | result: success | client_id: %v", c.config.ID)
		os.Exit(0)
    }()

	// There is an autoincremental msgID to identify every message sent
	// Messages if the message amount threshold has not been surpassed
	for msgID := 1; msgID <= c.config.LoopAmount; msgID++ {
		// Create the connection the server in every loop iteration. Send an
		c.createClientSocket()

		// Create Bet from env
		bet, err := from_env()
		if err != nil {
			fmt.Println("action: create_bet | result: fail | client_id: %v | msg: %v",
				c.config.ID,
				err,
			)
			return
		}

		// Send Bet to the server
		n, err := c.write_all(bet.to_bytes())
		if err != nil {
			log.Errorf("action: send_bet | result: fail | client_id: %v | msg: %v",
				c.config.ID,
				err,
			)
			return
		}
		log.Infof("action: apuesta_enviada | result: success | dni: %v | numero: %v",
			bet.dni,
			bet.bet,
		)

		if len(bet.ToBytes()) > n {
			log.Errorf("action: send_bet | result: fail | msg: short_read")
		}

		// Wait for response from server and close the connection
		msg, err := bufio.NewReader(c.conn).ReadString('\n')
		c.conn.Close()

		if err != nil {
			log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v",
				c.config.ID,
				err,
			)
			return
		}

		log.Infof("action: receive_message | result: success | client_id: %v | msg: %v",
			c.config.ID,
			msg,
		)

		// Wait a time between sending one message and the next one
		time.Sleep(c.config.LoopPeriod)

	}
	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
}

// To avoid 'short write'
func (c *Client) write_all(buffer []byte) (int, error) {
	bytes_sent := 0
	for bytes_sent < len(buffer) {
		n, err := c.conn.Write(buffer[bytes_sent:])
		if err != nil {
			return bytes_sent, err
		}
		bytes_sent += n
	}
	return bytes_sent, nil
}