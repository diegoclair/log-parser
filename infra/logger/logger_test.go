package logger

import (
	"testing"

	"github.com/diegoclair/log-parser/infra/config"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	cfg := &config.Config{}
	logger := New(cfg)
	require.NotNil(t, logger)
}
