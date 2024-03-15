package contract

import (
	"context"

	"github.com/diegoclair/log-parser/application/dto"
)

type Writer interface {
	StartWriting(ctx context.Context, data <-chan dto.Report)
}
