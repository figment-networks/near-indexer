package cli

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/olekukonko/tablewriter"

	"github.com/figment-networks/near-indexer/config"
	"github.com/figment-networks/near-indexer/near"
	"github.com/figment-networks/near-indexer/store"
)

func startStatus(cfg *config.Config) error {
	db, err := initStore(cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	rpc := near.DefaultClient(cfg.RPCEndpoint)

	status, err := rpc.Status()
	if err != nil {
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeader([]string{"Attribute", "Value"})

	table.AppendBulk([][]string{
		{"Chain ID", status.ChainID},
		{"Version", status.Version.String()},
		{"Syncing", fmt.Sprintf("%v", status.SyncInfo.Syncing)},
		{"Latest Block Height", fmt.Sprintf("%v", status.SyncInfo.LatestBlockHeight)},
		{"Latest Block Hash", status.SyncInfo.LatestBlockHash},
		{"Latest Block Time", status.SyncInfo.LatestBlockTime.UTC().Format(time.RFC3339)},
	})

	lastBlock, err := db.Blocks.Last()
	if err == nil {
		table.AppendBulk([][]string{
			{"Indexer Block Height", fmt.Sprintf("%v", lastBlock.Height)},
			{"Indexer Block Hash", lastBlock.Hash},
			{"Indexer Block Time", lastBlock.Time.UTC().Format(time.RFC3339)},
			{"Indexer Lag", fmt.Sprintf("%v", status.SyncInfo.LatestBlockHeight-uint64(lastBlock.Height))},
		})
	} else {
		if err == store.ErrNotFound {
			table.Append([]string{"Indexed Block", "N/A"})
		}
		log.Println("cant fetch recent block:", err)
	}

	genesis, err := rpc.GenesisConfig()
	if err == nil {
		table.AppendBulk([][]string{
			{"Genesis Block Height", fmt.Sprintf("%v", genesis.GenesisHeight)},
			{"Genesis Block Time", genesis.GenesisTime.UTC().Format(time.RFC3339)},
		})
	} else {
		log.Println("cant fetch genesis config:", err)
	}

	table.Render()
	return nil
}
