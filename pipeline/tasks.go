package pipeline

import (
	"time"

	"github.com/sirupsen/logrus"
)

const (
	FetcherTaskName   = "fetcher"
	ParserTaskName    = "parser"
	PersistorTaskName = "persistor"
	AnalyzerTaskName  = "analyzer"
)

func logTaskDuration(name string, ts time.Time) {
	logrus.
		WithField("task", name).
		WithField("duration", time.Since(ts).Milliseconds()).
		Debug("task finished")
}
