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

	rpc := near.NewClient(cfg.RPCEndpoint)
	rpc.SetDebug(cfg.Debug)

	height, err := store.Heights.LastSuccessful()
	if err != nil {
		return err
	}

	heightStatuses, err := store.Heights.StatusCounts()
	if err != nil {
		return err
	}

	fmt.Println("=== Height Indexing ===")
	fmt.Println("Last height:", height.Height)
	for _, s := range heightStatuses {
		fmt.Printf("Status: %s, Count: %d\n", s.Status, s.Num)
	}

	status, err := rpc.Status()
	if err != nil {
		return err
	}
	info := status.SyncInfo

	fmt.Println("=== Node Status ===")
	fmt.Println("Chain:", status.ChainID)
	fmt.Println("Syncing:", info.Syncing)
	fmt.Println("Last height:", status.SyncInfo.LatestBlockHeight)
	fmt.Println("Last hash:", status.SyncInfo.LatestBlockHash)
	fmt.Println("Last time:", status.SyncInfo.LatestBlockTime)

	gc, err := rpc.GenesisConfig()
	if err != nil {
		return err
	}
	fmt.Println("=== Genesis Status ===")
	fmt.Println("Height:", gc.GenesisHeight)
	fmt.Println("Time:", gc.GenesisTime)

	return nil
}
