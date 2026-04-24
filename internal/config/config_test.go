package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/driftwatch/driftwatch/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".driftwatch.json")
	require.NoError(t, os.WriteFile(path, []byte(content), 0o644))
	return path
}

func TestLoad_ValidConfig(t *testing.T) {
	path := writeTemp(t, `{
		"version": "1",
		"profiles": {
			"prod": {"name": "prod", "provider": "aws", "region": "us-east-1"}
		}
	}`)

	cfg, err := config.Load(path)
	require.NoError(t, err)
	assert.Equal(t, "1", cfg.Version)
	assert.Equal(t, "aws", cfg.Profiles["prod"].Provider)
	assert.Equal(t, "us-east-1", cfg.Profiles["prod"].Region)
}

func TestLoad_MissingVersion(t *testing.T) {
	path := writeTemp(t, `{
		"profiles": {
			"prod": {"name": "prod", "provider": "aws"}
		}
	}`)

	_, err := config.Load(path)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "version field is required")
}

func TestLoad_MissingProvider(t *testing.T) {
	path := writeTemp(t, `{
		"version": "1",
		"profiles": {
			"staging": {"name": "staging"}
		}
	}`)

	_, err := config.Load(path)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "provider is required")
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := config.Load("/nonexistent/path/.driftwatch.json")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "reading")
}

func TestLoad_InvalidJSON(t *testing.T) {
	path := writeTemp(t, `{not valid json}`)
	_, err := config.Load(path)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "parsing")
}
