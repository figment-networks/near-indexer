package cli

import (
	"flag"
	"fmt"
	"log"

	"github.com/figment-networks/near-indexer/config"
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
		log.Println(versionString())
		return
	}

	cfg, err := initConfig(configPath)
	if err != nil {
		terminate(err)
	}

	if runCommand == "" {
		terminate("Command is required")
	}

	if err := startCommand(cfg, runCommand); err != nil {
		terminate(err)
	}
}

func startCommand(cfg *config.Config, name string) error {
	switch name {
	case "migrate":
		return startMigrations(cfg)
	case "server":
		return startServer(cfg)
	case "worker":
		return startWorker(cfg)
	case "sync":
		return runSync(cfg)
	case "status":
		return startStatus(cfg)
	default:
		return fmt.Errorf("%s is not a valid command", name)
	}
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
