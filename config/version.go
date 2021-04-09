package config

import "fmt"

var (
	AppName    = "near-indexer"
	AppVersion = "0.3.4"
	GitCommit  = "-"
	GoVersion  = "-"
)

// VersionString returns the full app version string
func VersionString() string {
	return fmt.Sprintf(
		"%s %s (git: %s, %s)",
		AppName,
		AppVersion,
		GitCommit,
		GoVersion,
	)
}
