package naming

import (
	"errors"
	"github.com/conplementag/cops-hq/v2/internal"
	"github.com/conplementag/cops-hq/v2/pkg/naming/patterns"
)

// New creates a new naming convention service.
// Parameters 'context', 'color', 'region' and 'environment' are mandatory.
//
//	Context should be set to the application name. In case of a complex system with multiple modules / subsystems,
//		'module' should be set as well (but it is not required).
//	Color can be used for blue/green infrastructure deployments, useful in disaster recovery scenarios. For example, if
//		you would ever need to recreate infrastructure itself and perform a disaster recovery, you would get naming clashes
//		for globally-unique resources. When you set the color property from the beginning, for example to 'b' (as in blue), you
//		could change it in the future for 'g' (as in green) to redeploy everything. Keep the color names very short due to limits
//		in Azure for the resource name lenghts.
//	Environment provides you the possibility to isolate your resources per environment, e.g. prod, int, stage, dev, etc.
//	Regions should be in form of Azure long regions string, e.g. westeurope or northeurope.
//
// Per default, normal naming convention pattern is used. If required, you can override the pattern using the
// SetPattern() method
func New(context string, module string, color string, region string, environment string) (*Service, error) {
	if context == "" {
		return nil, internal.ReturnErrorOrPanic(errors.New("context must be provided"))
	}

	if color == "" {
		return nil, internal.ReturnErrorOrPanic(errors.New("color must be provided, recommended: use either b for blue or g for green"))
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
		color:       color,
		region:      region,
		environment: environment,
	}, nil
}
