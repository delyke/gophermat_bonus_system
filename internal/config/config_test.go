package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/delyke/gophermat_bonus_system/internal/config/defaults"
)

func TestLoadDefaults(t *testing.T) {
	originalArgs := os.Args
	os.Args = []string{os.Args[0]}
	t.Cleanup(func() {
		os.Args = originalArgs
	})

	err := Load()
	require.NoError(t, err)

	cfg := Get()
	require.Equal(t, defaults.JWTSecret, cfg.JWT.Secret())
	require.Equal(t, 24*time.Hour, cfg.JWT.TTL())
}

func TestLoadEnvOverrides(t *testing.T) {
	originalArgs := os.Args
	os.Args = []string{os.Args[0]}
	t.Cleanup(func() {
		os.Args = originalArgs
	})

	t.Setenv("JWT_SECRET", "custom-secret")
	t.Setenv("JWT_TTL", "2h")

	err := Load()
	require.NoError(t, err)

	cfg := Get()
	require.Equal(t, "custom-secret", cfg.JWT.Secret())
	require.Equal(t, 2*time.Hour, cfg.JWT.TTL())
}
