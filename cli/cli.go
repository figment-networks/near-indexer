package cli

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/sirupsen/logrus"

	"github.com/figment-networks/near-indexer/config"
	"github.com/figment-networks/near-indexer/near"
	"github.com/figment-networks/near-indexer/server"
	"github.com/figment-networks/near-indexer/store"
)

// Run executes the command line interface
func Run() {
	var (
		configPath  string
		runCommand  string
		showVersion bool
	)
	runCommand = "server"
	configPath = "./config.example.json"
	flag.Parse()

	if showVersion {
		log.Println(config.VersionString())
		return
	}

	cfg, err := initConfig(configPath)
	if err != nil {
		terminate(err)
	}

	logger := initLogger(cfg)

	config.InitRollbar(cfg)
	defer config.TrackRecovery()

	if runCommand == "" {
		terminate("Command is required")
	}

	if cfg.Debug {
		initProfiler()
	}

	if err := startCommand(cfg, logger, runCommand); err != nil {
		terminate(err)
	}
}

func startCommand(cfg *config.Config, logger *logrus.Logger, name string) error {
	switch name {
	case "server":
		return startServer(cfg, logger)
	case "worker":
		return startWorker(cfg, logger)
	case "sync":
		return runSync(cfg, logger)
	case "status":
		return startStatus(cfg)
	case "cleanup":
		return startCleanup(cfg, logger)
	case "reset":
		return startReset(cfg)
	default:
		return fmt.Errorf("%s is not a valid command", name)
	}
}

func confirm(message string) bool {
	var val string

	fmt.Printf("%s (Y/N): ", message)
	fmt.Fscanln(os.Stdin, &val)

	return strings.TrimSpace(val) == "Y"
}

func terminate(message interface{}) {
	if message != nil {
		log.Fatal("ERROR: ", message)
	}
}

func initConfig(path string) (*config.Config, error) {
	cfg := config.New()

	if err := config.FromEnv(cfg); err != nil {
		return nil, err
	}

	if path != "" {
		if err := config.FromFile(path, cfg); err != nil {
			return nil, err
		}
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func initLogger(cfg *config.Config) *logrus.Logger {
	logger := logrus.StandardLogger()

	switch cfg.LogLevel {
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	}

	return logger
}

func initClient(cfg *config.Config) near.Client {
	client := near.DefaultClient(cfg.RPCEndpoint)
	client.SetTimeout(cfg.RPCClientTimeout())
	if cfg.LogLevel == "debug" {
		client.SetDebug(true)
	}
	return client
}

func initStore(cfg *config.Config) (*store.Store, error) {
	db, err := store.New(cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	if cfg.LogLevel == "debug" {
		db.SetDebugMode(true)
	}

	return db, nil
}

func initSignals() chan os.Signal {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)
	return c
}

func initProfiler() {
	go server.StartProfiler()
}

func initRollbar(cfg *config.Config) {
	config.InitRollbar(cfg)
}
