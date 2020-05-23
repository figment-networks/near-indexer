package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/kelseyhightower/envconfig"
)

const (
	modeDevelopment = "development"
	modeProduction  = "production"
)

var (
	errEndpointRequired        = errors.New("Near RPC endpoint is required")
	errDatabaseRequired        = errors.New("Database credentials are required")
	errSyncIntervalRequired    = errors.New("Sync interval is required")
	errSyncIntervalInvalid     = errors.New("Sync interval is invalid")
	errCleanupIntervalRequired = errors.New("Cleanup interval is required")
	errCleanupIntervalInvalid  = errors.New("Cleanup interval is invalid")
)

// Config holds the configration data
type Config struct {
	AppEnv           string `json:"app_env" envconfig:"APP_ENV" default:"development"`
	RPCEndpoint      string `json:"rpc_endpoint" envconfig:"RPC_ENDPOINT"`
	ServerAddr       string `json:"server_addr" envconfig:"SERVER_ADDR" default:"0.0.0.0"`
	ServerPort       int    `json:"server_port" envconfig:"SERVER_PORT" default:"8081"`
	SyncInterval     string `json:"sync_interval" envconfig:"SYNC_INTERVAL" default:"10s"`
	CleanupInterval  string `json:"cleanup_interval" envconfig:"CLEANUP_INTERVAL" default:"10m"`
	CleanupThreshold int    `json:"cleanup_threshold" envconfig:"CLEANUP_THRESHOLD" default:"1000"`
	DatabaseURL      string `json:"database_url" envconfig:"DATABASE_URL"`
	DumpDir          string `json:"dump_dir" envconfig:"DUMP_DIR"`
	Debug            bool   `json:"debug" envconfig:"DEBUG"`
	LogLevel         string `json:"log_level" envconfig:"LOG_LEVEL" default:"info"`

	syncDuration    time.Duration
	cleanupDuration time.Duration
}

// Validate returns an error if config is invalid
func (c *Config) Validate() error {
	if c.RPCEndpoint == "" {
		return errEndpointRequired
	}

	if c.DatabaseURL == "" {
		return errDatabaseRequired
	}

	if c.SyncInterval == "" {
		return errSyncIntervalRequired
	}
	if d, err := time.ParseDuration(c.SyncInterval); err != nil {
		return errSyncIntervalInvalid
	} else {
		c.syncDuration = d
	}

	if c.CleanupInterval == "" {
		return errCleanupIntervalRequired
	}
	if d, err := time.ParseDuration(c.CleanupInterval); err != nil {
		return errCleanupIntervalInvalid
	} else {
		c.cleanupDuration = d
	}

	return nil
}

// IsDevelopment returns true if app is in dev mode
func (c *Config) IsDevelopment() bool {
	return c.AppEnv == modeDevelopment
}

// IsProduction returns true if app is in production mode
func (c *Config) IsProduction() bool {
	return c.AppEnv == modeProduction
}

// ListenAddr returns a full listen address and port
func (c *Config) ListenAddr() string {
	return fmt.Sprintf("%s:%d", c.ServerAddr, c.ServerPort)
}

func (c *Config) SyncDuration() time.Duration {
	return c.syncDuration
}

func (c *Config) CleanupDuration() time.Duration {
	return c.cleanupDuration
}

// New returns a new config
func New() *Config {
	return &Config{}
}

// FromFile reads the config from a file
func FromFile(path string, config *Config) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, config)
}

// FromEnv reads the config from environment variables
func FromEnv(config *Config) error {
	return envconfig.Process("", config)
}
