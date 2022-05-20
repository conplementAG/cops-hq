package naming

import (
	"errors"
	"github.com/conplementag/cops-hq/internal"
	"github.com/conplementag/cops-hq/pkg/naming/patterns"
)

// New creates a new naming convention service.
// Parameters 'context', 'region' and 'environment' are mandatory.
//     Context should be set to the application name. In case of a complex system with multiple modules / subsystems,
//       'module' should be set as well.
//     Environment provides you the possibility to isolate your resources per environment, e.g. prod, int, stage, dev, etc.
//     Regions should be in form of Azure long regions string, e.g. westeurope or northeurope.
// Per default, normal naming convention pattern is used. If required, you can override the pattern using the
// SetPattern() method
func New(context string, region string, environment string, module string) (*Service, error) {
	if context == "" {
		return nil, internal.ReturnErrorOrPanic(errors.New("context must be provided"))
	}

	if region == "" {
		return nil, internal.ReturnErrorOrPanic(errors.New("region must be provided"))
	}

	if environment == "" {
		return nil, internal.ReturnErrorOrPanic(errors.New("environment must be provided"))
	}

	return &Service{
		pattern:     patterns.Normal,
		context:     context,
		module:      module,
		region:      region,
		environment: environment,
	}, nil
}
