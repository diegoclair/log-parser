package service

import (
	"github.com/diegoclair/go_utils/logger"
	"github.com/diegoclair/log-parser/application/contract"
)

type Services struct {
	QuakeService contract.QuakeService
}

type service struct {
	log logger.Logger
}

// New to get instance of all services
func New(log logger.Logger) (*Services, error) {
	svc := &service{
		log: log,
	}

	return &Services{
		QuakeService: newQuakeService(svc),
	}, nil
}
