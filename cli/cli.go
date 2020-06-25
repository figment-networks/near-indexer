package cli

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/figment-networks/near-indexer/config"
	"github.com/figment-networks/near-indexer/server"
	"github.com/figment-networks/near-indexer/store"
)

// Run executes the command line interface
func Run() {
	var configPath string
	var runCommand string
	var showVersion bool

	flag.BoolVar(&showVersion, "v", false, "Show application version")
	flag.StringVar(&configPath, "config", "", "Path to config")
	flag.StringVar(&runCommand, "cmd", "", "Command to run")
	flag.Parse()

	if showVersion {
		log.Println(config.VersionString())
		return
	}

	cfg, err := initConfig(configPath)
	if err != nil {
		terminate(err)
	}

	config.InitRollbar(cfg)
	defer config.TrackRecovery()

	if runCommand == "" {
		terminate("Command is required")
	}

	if cfg.Debug {
		initProfiler()
	}

	if err := startCommand(cfg, runCommand); err != nil {
		terminate(err)
	}
}

func startCommand(cfg *config.Config, name string) error {
	switch name {
	case "migrate", "migrate:up", "migrate:down", "migrate:redo":
		return startMigrations(name, cfg)
	case "server":
		return startServer(cfg)
	case "worker":
		return startWorker(cfg)
	case "sync":
		return runSync(cfg)
	case "status":
		return startStatus(cfg)
	case "cleanup":
		return startCleanup(cfg)
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

func initStore(cfg *config.Config) (*store.Store, error) {
	db, err := store.New(cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	db.SetDebugMode(cfg.Debug)

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
