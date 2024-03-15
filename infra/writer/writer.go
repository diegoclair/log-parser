package writer

import (
	"context"
	"encoding/json"
	"io"

	"github.com/diegoclair/go_utils/logger"
	"github.com/diegoclair/log-parser/application/contract"
	"github.com/diegoclair/log-parser/application/dto"
)

type writer struct {
	file io.Writer
	log  logger.Logger
}

func NewWriter(ctx context.Context, file io.Writer, log logger.Logger) contract.Writer {
	return &writer{
		file: file,
		log:  log,
	}
}

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
