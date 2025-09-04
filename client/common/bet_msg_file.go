package common

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

const BET_FIELDS = 5

type BetLoader struct {
	reader   			*csv.Reader
	file     			*os.File
	agency_id 			string
	max_bets_per_batch 	int
}

// Load bets from CSV files using our BetLoader
func LoadBets(max_bets_per_batch int) (*BetLoader, error) {
	agency_id := os.Getenv("CLI_ID")
	if agency_id == "" {
		return nil, fmt.Errorf("Please define CLI_ID env variable")
	}

	csv_path := fmt.Sprintf("./agency-%s.csv", agency_id)
	file, err := os.Open(csv_path)
	if err != nil {
		return nil, fmt.Errorf("Error opening file: %v", err)
	}

	reader := csv.NewReader(file)
	return &BetLoader{
		reader:    			reader,
		file:      			file,
		agency_id: 			agency_id,
		max_bets_per_batch: max_bets_per_batch,
	}, nil
}

// Read 135 lines (batch) and return a byte buffer and it's total length
// Reads from the last point read
func (br *BetLoader) LoadBetBatch() ([]byte, uint32, bool, error) {
	var buffer bytes.Buffer
	eof := false
	lineCount := 0

	for lineCount < br.max_bets_per_batch {
		record, err := br.reader.Read()
		if err == io.EOF {
			eof = true
			break
		}
		if err != nil {
			return nil, 0, eof, fmt.Errorf("Error reading file: %v", err)
		}
		if len(record) < BET_FIELDS {
			return nil, 0, eof, fmt.Errorf("Invalid format on line %d", lineCount+1)
		}

		bet := Bet{
			id_agency:   br.agencyID,
			name:  record[0],
			surname:   record[1],
			dni: record[2],
			birth_date:  record[3],
			bet:  record[4],
		}

		buffer.Write(bet.to_bytes())
		lineCount++
	}

	return buffer.Bytes(), uint32(buffer.Len()), eof, nil
}

// Close file after reading all the batches.
func (br *BetLoader) Close() {
	br.file.Close()
}