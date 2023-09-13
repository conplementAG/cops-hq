package sops

import "github.com/conplementag/cops-hq/v2/pkg/commands"

// New creates a new instance of Sops, which is a wrapper around common Sops functionality.
// Parameters:
//
//	executor (can be provided from hq.GetExecutor() or by instantiating your own),
func New(executor commands.Executor) Sops {

	return &sopsWrapper{
		executor: executor,
	}
}
