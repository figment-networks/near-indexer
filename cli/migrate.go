package cli

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pressly/goose"

	"github.com/figment-networks/near-indexer/config"
	"github.com/figment-networks/near-indexer/migrations"
)

func startMigrations(cmd string, cfg *config.Config) error {
	store, err := initStore(cfg)
	if err != nil {
		return err
	}
	defer store.Close()

	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	for path, f := range migrations.Assets.Files {
		if filepath.Ext(path) != ".sql" {
			continue
		}

		extPath := filepath.Join(tmpDir, filepath.Base(path))
		if err := ioutil.WriteFile(extPath, f.Data, 0755); err != nil {
			return err
		}
	}

	dir := "up"
	if chunks := strings.Split(cmd, ":"); len(chunks) > 1 {
		dir = chunks[1]
	}

	switch dir {
	case "migrate", "up":
		err = goose.Up(store.Conn(), tmpDir)
	case "down":
		err = goose.Down(store.Conn(), tmpDir)
	case "redo":
		if err = goose.Down(store.Conn(), tmpDir); err != nil {
			return err
		}
		err = goose.Up(store.Conn(), tmpDir)
	}

	return err
}
