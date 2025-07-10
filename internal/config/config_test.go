package config

import (
	"os"
	"path/filepath"
	"testing"
)

const testConfigYAML = `
env: "test"
server:
  host: localhost
  port: "8080"
grpc:
  host: localhost
  port: "50051"
`

func TestMustLoadConfig_Success(t *testing.T) {
	// Создаём временный файл с конфигом
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	if err := os.WriteFile(configPath, []byte(testConfigYAML), 0644); err != nil {
		t.Fatalf("failed to write temp config: %v", err)
	}

	// Устанавливаем переменную окружения
	os.Setenv("CONFIG_PATH", configPath)
	defer os.Unsetenv("CONFIG_PATH")

	cfg := MustLoadConfig()

	if cfg.Env != "test" {
		t.Errorf("expected env 'test', got '%s'", cfg.Env)
	}
	if cfg.Server.Host != "localhost" {
		t.Errorf("expected host 'localhost', got '%s'", cfg.Server.Host)
	}
	if cfg.Server.Port != "8080" {
		t.Errorf("expected port '8080', got '%s'", cfg.Server.Port)
	}
}

func TestMustLoadConfig_EmptyPath(t *testing.T) {
	os.Unsetenv("CONFIG_PATH")
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic on empty config path")
		}
	}()
	MustLoadConfig()
}

func TestMustLoadConfig_FileNotExist(t *testing.T) {
	os.Setenv("CONFIG_PATH", "/non/existent/path.yaml")
	defer os.Unsetenv("CONFIG_PATH")
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic on missing config file")
		}
	}()
	MustLoadConfig()
}
