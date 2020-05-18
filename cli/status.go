package cli

import (
	"fmt"

	"github.com/figment-networks/near-indexer/config"
	"github.com/figment-networks/near-indexer/near"
)

func startStatus(cfg *config.Config) error {
	status, err := near.NewClient(cfg.RPCEndpoint).Status()
	if err != nil {
		terminate(err)
	}
	info := status.SyncInfo

	fmt.Println("Chain:", status.ChainID)
	fmt.Println("Syncing:", info.Syncing)
	fmt.Println("Last height:", status.SyncInfo.LatestBlockHeight)
	fmt.Println("Last hash:", status.SyncInfo.LatestBlockHash)
	fmt.Println("Last time:", status.SyncInfo.LatestBlockTime)

	return nil
}
