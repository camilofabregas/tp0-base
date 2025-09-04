package common

import (
	"bufio"
	"encoding/binary"
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
	ID            	string
	ServerAddress 	string
	LoopAmount    	int
	LoopPeriod    	time.Duration
	MaxBetsPerBatch int
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

	betLoader, err := LoadBets(c.config.MaxBetsPerBatch)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer betLoader.Close()

	go func() {
        <-sigc
        log.Infof("action: sigterm_received | result: in_progress | client_id: %v | msg: Starting client shutdown", c.config.ID)
        if c.conn != nil {
            c.conn.Close()
        }
		betLoader.Close()

		log.Infof("action: shutdown | result: success | client_id: %v", c.config.ID)
		os.Exit(0)
    }()

	for {
		// Create the connection the server in every loop iteration. Send an
		c.createClientSocket()
		data, bytes_total, eof, err := betLoader.LoadBetBatch()
		if err != nil {
			log.Errorf("action: batch_read | result: fail")
			return
		}
		log.Infof("action: batch_read | result: success")

		len_bytes := make([]byte, 4) // 4 bytes (32 bits)
		binary.BigEndian.PutUint32(len_bytes, bytes_total)
		batch := append(len_bytes, data...)

		// Send Bet Batch to the server
		n, err := c.write_all(batch)
		if err != nil {
			log.Errorf("action: batch_enviado | result: fail | client_id: %v | msg: %v",
				c.config.ID,
				err,
			)
			return
		}
		log.Infof("action: batch_enviado | result: success")

		if len(batch) > n {
			log.Errorf("action: batch_enviado | result: fail | msg: short_write")
		}

		// Wait for response from server and close the connection
		msg, err := bufio.NewReader(c.conn).ReadString('\n')

		if err != nil {
			log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v",
				c.config.ID,
				err,
			)
			return
		}

		if msg == "ACK BATCH\n" {
			log.Infof("action: batch_acknowledged | result: success | client_id: %v | msg: %v",
				c.config.ID,
				string(msg),
			)
		} else {
			log.Errorf("action: batch_acknowledged | result: fail | client_id: %v | msg: %v",
				c.config.ID,
				string(msg),
			)
		}

		if eof {
			log.Infof("action: envio_completado | result: success")
			c.conn.Close()
			break
		}
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