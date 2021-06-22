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
	modeStaging     = "staging"
	modeProduction  = "production"
)

var (
	errEndpointRequired        = errors.New("Near RPC endpoint is required")
	errDatabaseRequired        = errors.New("Database credentials are required")
	errSyncIntervalRequired    = errors.New("Sync interval is required")
	errSyncIntervalInvalid     = errors.New("Sync interval is invalid")
	errCleanupIntervalRequired = errors.New("Cleanup interval is required")
	errCleanupIntervalInvalid  = errors.New("Cleanup interval is invalid")
	errRPCTimeoutInvalid       = errors.New("RPC timeout interval is invalid")
)

// Config holds the configration data
type Config struct {
	AppEnv           string `json:"app_env" envconfig:"APP_ENV" default:"development"`
	RPCEndpoints     string `json:"rpc_endpoints" envconfig:"NEAR_RPC_ENDPOINTS"`
	RPCTimeout       string `json:"rpc_timeout" envconfig:"NEAR_RPC_TIMEOUT" default:"75s"`
	ServerAddr       string `json:"server_addr" envconfig:"SERVER_ADDR" default:"0.0.0.0"`
	ServerPort       int    `json:"server_port" envconfig:"SERVER_PORT" default:"8081"`
	StartHeight      uint64 `json:"start_height" envconfig:"START_HEIGHT"`
	SyncInterval     string `json:"sync_interval" envconfig:"SYNC_INTERVAL" default:"500ms"`
	SyncBatchSize    int    `json:"sync_batch_size" envconfig:"SYNC_BATCH_SIZE" default:"10"`
	CleanupInterval  string `json:"cleanup_interval" envconfig:"CLEANUP_INTERVAL" default:"10m"`
	CleanupThreshold int    `json:"cleanup_threshold" envconfig:"CLEANUP_THRESHOLD" default:"3600"`
	DatabaseURL      string `json:"database_url" envconfig:"DATABASE_URL"`
	DumpDir          string `json:"dump_dir" envconfig:"DUMP_DIR"`
	Debug            bool   `json:"debug" envconfig:"DEBUG"`
	LogLevel         string `json:"log_level" envconfig:"LOG_LEVEL" default:"info"`

	// delegation calls
	RetryCountDlg    int `json:"retry_count_delegation_calls" envconfig:"RETRY_COUNT_DELEGATION_CALLS" default:"4"`
	ConcurrencyLevel int `json:"concurrency_level" envconfig:"CONCURRENCY_LEVEL" default:"2"`

	// Exception tracking
	RollbarToken     string `json:"rollbar_token" envconfig:"ROLLBAR_TOKEN"`
	RollbarNamespace string `json:"rollbar_namespace" envconfig:"ROLLBAR_NAMESPACE"`

	syncDuration    time.Duration
	cleanupDuration time.Duration
	rpcTimeout      time.Duration
}

// Validate returns an error if config is invalid
func (c *Config) Validate() error {
	if c.RPCEndpoints == "" {
		return errEndpointRequired
	}

	if c.DatabaseURL == "" {
		return errDatabaseRequired
	}

	if c.SyncInterval == "" {
		return errSyncIntervalRequired
	}

	d, err := time.ParseDuration(c.SyncInterval)
	if err != nil {
		return errSyncIntervalInvalid
	}
	c.syncDuration = d

	if c.CleanupInterval == "" {
		return errCleanupIntervalRequired
	}
	d, err = time.ParseDuration(c.CleanupInterval)
	if err != nil {
		return errCleanupIntervalInvalid
	}
	c.cleanupDuration = d

	rpcTimeout, err := time.ParseDuration(c.RPCTimeout)
	if err != nil {
		return errRPCTimeoutInvalid
	}
	c.rpcTimeout = rpcTimeout

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

// SyncDuration returns the parsed duration for the sync pipeline
func (c *Config) SyncDuration() time.Duration {
	return c.syncDuration
}

// CleanupDuration returns the parsed duration for the cleanup pipeline
func (c *Config) CleanupDuration() time.Duration {
	return c.cleanupDuration
}

// RPCClientTimeout returns the timeout value for RPC calls
func (c *Config) RPCClientTimeout() time.Duration {
	return c.rpcTimeout
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
