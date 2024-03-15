package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetDefaultConfig(t *testing.T) {
	cfg := GetDefaultConfig()
	require.NotNil(t, cfg)
	require.Equal(t, "log-parser", cfg.AppName)
	require.True(t, cfg.LogDebug)
}
