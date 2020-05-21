package cli

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pressly/goose"

	"github.com/figment-networks/near-indexer/config"
	"github.com/figment-networks/near-indexer/migrations"
)

func startMigrations(cfg *config.Config) error {
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

	return goose.Up(store.Conn(), tmpDir)
}
