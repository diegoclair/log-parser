package main

import (
	"context"
	"flag"
	"os"
	"sync"
	"time"

	"github.com/diegoclair/log-parser/application/dto"
	"github.com/diegoclair/log-parser/application/service"
	"github.com/diegoclair/log-parser/infra/config"
	"github.com/diegoclair/log-parser/infra/logger"
	"github.com/diegoclair/log-parser/infra/writer"
	"github.com/diegoclair/log-parser/transport/scripts"
)

var (
	logPath string
)

func init() {
	flag.StringVar(&logPath, "logpath", "./qgames.log", "Quake log file path")
}

func main() {
	flag.Parse()

	start := time.Now()

	cfg := config.GetDefaultConfig()

	ctx := context.Background()
	log := logger.New(cfg)

	logFile, err := os.Open(logPath)
	if err != nil {
		log.Errorf(ctx, "Error to open file: %v", err)
		return
	}

	defer logFile.Close()

	resultFile, err := os.Create("./result.json")
	if err != nil {
		log.Errorf(ctx, "Error to create file: %v", err)
	}

	defer resultFile.Close()

	rawLinesChan := make(chan string)
	writerChan := make(chan dto.Report)

	svc, err := service.New(log)
	if err != nil {
		log.Errorf(ctx, "Error getting NewAuthToken: %v", err)
		return
	}

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		writer.NewWriter(resultFile, log).StartWriting(ctx, writerChan)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		svc.QuakeService.StartExtractingData(ctx, rawLinesChan, writerChan)
	}()

	scripts.NewQuakeLogParser(log).ReadLinesFromQuakeLog(ctx, logFile, rawLinesChan)

	wg.Wait()

	log.Infof(ctx, "Execution time: %v", time.Since(start))
}
