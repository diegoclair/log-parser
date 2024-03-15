package scripts

import (
	"bufio"
	"context"
	"os"

	"github.com/diegoclair/go_utils/logger"
)

type QuakeLogParser struct {
	log logger.Logger
}

func NewQuakeLogParser(log logger.Logger) *QuakeLogParser {
	return &QuakeLogParser{
		log: log,
	}
}

func (q *QuakeLogParser) ExecuteQuakeLogParser(ctx context.Context, lineChan chan<- string) {
	file, err := os.Open("./qgames.log")
	if err != nil {
		q.log.Errorf(ctx, "Error to open file: %v", err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lineChan <- scanner.Text()
	}

	close(lineChan)
}
