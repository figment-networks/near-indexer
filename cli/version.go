package cli

import "fmt"

var (
	appName    = "near-indexer"
	appVersion = "0.1.0"
	gitCommit  = "-"
	goVersion  = "-"
)

func versionString() string {
	return fmt.Sprintf(
		"%s %s (git: %s, %s)",
		appName,
		appVersion,
		gitCommit,
		goVersion,
	)
}
