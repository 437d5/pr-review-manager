package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMustLoadConfig_WithEnvOverrides(t *testing.T) {
	t.Setenv("MODE", "prod")
	t.Setenv("REVIEWER_ADDRESS", ":9090")
	t.Setenv("REVIEWER_READ_TIMEOUT", "30")
	t.Setenv("REVIEWER_WRITE_TIMEOUT", "25")
	t.Setenv("REVIEWER_IDLE_TIMEOUT", "120")
	t.Setenv("DB_NAME", "custom_db")
	t.Setenv("DB_HOST", "db-host")
	t.Setenv("DB_PORT", "6543")
	t.Setenv("DB_USER", "db-user")
	t.Setenv("DB_PASSWORD", "db-pass")

	cfg := MustLoadConfig()

	require.NotNil(t, cfg)
	assert.Equal(t, "prod", cfg.Mode)
	assert.Equal(t, ":9090", cfg.Address)
	assert.Equal(t, 30, cfg.ReadTimeout)
	assert.Equal(t, 25, cfg.WriteTimeout)
	assert.Equal(t, 120, cfg.IdleTimeout)

	assert.Equal(t, "custom_db", cfg.DB.Name)
	assert.Equal(t, "db-host", cfg.DB.Host)
	assert.Equal(t, "6543", cfg.DB.Port)
	assert.Equal(t, "db-user", cfg.DB.User)
	assert.Equal(t, "db-pass", cfg.DB.Password)
}

func TestMustLoadConfig_InvalidEnvPanics(t *testing.T) {
	t.Setenv("REVIEWER_READ_TIMEOUT", "invalid")

	defer t.Setenv("REVIEWER_READ_TIMEOUT", "")

	assert.Panics(t, func() {
		MustLoadConfig()
	})
}

func TestConfig_GetConnectionString(t *testing.T) {
	cfg := &Config{
		DB: DBConfig{
			User:     "postgres",
			Password: "secret",
			Host:     "localhost",
			Port:     "5433",
			Name:     "test_db",
		},
	}

	dsn := cfg.GetConnectionString()

	assert.Equal(t, "postgres://postgres:secret@localhost:5433/test_db?sslmode=disable", dsn)
}
