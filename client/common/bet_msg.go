package common

import (
	"encoding/binary"
	"fmt"
	"os"
)

// Struct of the bet that the client will send to the server
type Bet struct {
	id_agency    string
	name         string
	surname      string
	dni         string
	birth_date   string
	bet   string
}

// Create Bet from environment variables
func from_env() (*Bet, error) {
	bet := &Bet{
		id_agency:    os.Getenv("CLI_ID"),
		name:         os.Getenv("NOMBRE"),
		surname:      os.Getenv("APELLIDO"),
		dni:          os.Getenv("DOCUMENTO"),
		birth_date:   os.Getenv("NACIMIENTO"),
		bet:          os.Getenv("NUMERO"),
	}

	if bet.id_agency == "" || bet.name == "" || bet.surname == "" ||
		bet.dni == "" || bet.birth_date == "" || bet.bet == "" {
		return nil, fmt.Errorf("Please provide all env variables")
	}

	return bet, nil
}

// Convert Bet to bytes for sending
func (b *Bet) to_bytes() []byte {
	msg := []byte(fmt.Sprintf("%s|%s|%s|%s|%s|%s\n", b.id_agency, b.name, b.surname, b.dni, b.birth_date, b.bet))
	len := uint32(len(msg))
	len_bytes := make([]byte, 4)
	binary.BigEndian.PutUint32(len_bytes, len)
	result := append(len_bytes, msg...)
	return result
}