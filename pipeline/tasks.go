package pipeline

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	fetcherTaskName   = "fetcher"
	parserTaskName    = "parser"
	persistorTaskName = "persistor"
	analyzerTaskName  = "analyzer"
)

type Task interface {
	Name() string
	ShouldRun(*Payload) bool
	Run(context.Context, *Payload) error
}

func logTaskDuration(t Task, ts time.Time) {
	logrus.
		WithField("task", t.Name()).
		WithField("duration", time.Since(ts).Milliseconds()).
		Debug("task finished")
}
