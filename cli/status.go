package cli

import (
	"fmt"

	"github.com/figment-networks/near-indexer/config"
	"github.com/figment-networks/near-indexer/near"
)

func startStatus(cfg *config.Config) error {
	store, err := initStore(cfg)
	if err != nil {
		return err
	}
	defer store.Close()

	heightStatuses, err := store.Heights.StatusCounts()
	if err != nil {
		return err
	}

	fmt.Println("=== Height Indexing ===")
	for _, s := range heightStatuses {
		fmt.Printf("Status: %s, Count: %d\n", s.Status, s.Num)
	}

	status, err := near.NewClient(cfg.RPCEndpoint).Status()
	if err != nil {
		terminate(err)
	}
	info := status.SyncInfo

	fmt.Println("=== Node Status ===")
	fmt.Println("Chain:", status.ChainID)
	fmt.Println("Syncing:", info.Syncing)
	fmt.Println("Last height:", status.SyncInfo.LatestBlockHeight)
	fmt.Println("Last hash:", status.SyncInfo.LatestBlockHash)
	fmt.Println("Last time:", status.SyncInfo.LatestBlockTime)

	return nil
}
