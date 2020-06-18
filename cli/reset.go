package cli

import (
	"errors"

	"github.com/figment-networks/near-indexer/config"
)

func startReset(cfg *config.Config) error {
	if !confirm("Are you sure you want to reset data?") {
		return errors.New("aborted")
	}

	db, err := initStore(cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	return db.ResetAll()
}
