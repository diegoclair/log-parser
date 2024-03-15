package scripts

import (
	"bufio"
	"context"
	"io"

	"github.com/diegoclair/go_utils/logger"
)

// QuakeLogParser is a struct that represents a Quake log parser.
type QuakeLogParser struct {
	log logger.Logger
}

// NewQuakeLogParser creates a new instance of QuakeLogParser.
func NewQuakeLogParser(log logger.Logger) *QuakeLogParser {
	return &QuakeLogParser{
		log: log,
	}
}

// ReadLinesFromQuakeLog reads lines from a Quake log file and sends them to lineChan channel.
func (q *QuakeLogParser) ReadLinesFromQuakeLog(ctx context.Context, file io.Reader, lineChan chan<- string) {
	defer close(lineChan)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lineChan <- scanner.Text()
	}
}
