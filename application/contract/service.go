package contract

import (
	"context"

	"github.com/diegoclair/log-parser/application/dto"
)

type QuakeService interface {
	// StartExtractingData to start extracting data from log line received from channel and create a report to be sent to writer
	StartExtractingData(ctx context.Context, rawLinesChan <-chan string, writerChan chan<- dto.Report)
}
