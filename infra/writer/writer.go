package writer

import (
	"context"
	"encoding/json"
	"io"

	"github.com/diegoclair/go_utils/logger"
	"github.com/diegoclair/log-parser/application/contract"
	"github.com/diegoclair/log-parser/application/dto"
)

// writer is a struct that implements the contract.Writer interface.
type writer struct {
	file io.Writer
	log  logger.Logger
}

// NewWriter creates a new instance of writer that implements contract.Writer interface.
func NewWriter(file io.Writer, log logger.Logger) contract.Writer {
	return &writer{
		file: file,
		log:  log,
	}
}

// StartWriting starts the process of writing the data received from the data channel to the specified file.
// It marshals the received data into JSON format and writes it to the file.
// If there is an error during marshaling or writing, it logs the error using the logger.Logger.
func (w *writer) StartWriting(ctx context.Context, data <-chan dto.Report) {
	reports := make(dto.QuakeDataReport)
	for report := range data {
		reports[report.GameName] = report
	}

	jsonData, err := json.Marshal(reports)
	if err != nil {
		w.log.Errorf(ctx, "Error to marshal data: %v", err)
	}

	_, err = w.file.Write(jsonData)
	if err != nil {
		w.log.Errorf(ctx, "Error to write data: %v", err)
	}
}
