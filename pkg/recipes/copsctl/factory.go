package copsctl

import (
	"github.com/conplementag/cops-hq/pkg/commands"
)

// New creates a new Copsctl recipe instance. Required parameters are:
//     executor (can be provided from hq.GetExecutor() or by instantiating your own)
func New(executor commands.Executor) Copsctl {
	return &copsctl{
		executor: executor,
	}
}
