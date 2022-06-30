package hq

import "errors"

type HqOptions struct {
	// Quiet HQ will create a quiet executor, piping all commands and outputs to the log file, but the console will be
	// kept clean. If needed, like in CI, a viper flag "verbose" can be used to override this behavior.
	Quiet bool

	// Filename for the log file, if enabled
	LogFileName string

	// Default logging to file can be disabled with this flag
	DisableFileLogging bool
}

func (options *HqOptions) Validate() error {
	if options.LogFileName == "" && !options.DisableFileLogging {
		return errors.New("you need to define the logFileName if logging to the file is enabled")
	}

	return nil
}
