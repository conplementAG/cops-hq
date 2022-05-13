package naming

import (
	"errors"
	"github.com/ahmetb/go-linq"
	"github.com/conplementag/cops-hq/pkg/naming/patterns"
	"github.com/conplementag/cops-hq/pkg/naming/regions"
	"github.com/conplementag/cops-hq/pkg/naming/resources"
	"strings"
)

type Service struct {
	pattern     patterns.Pattern
	context     string
	module      string
	region      string
	environment string
}

// New creates a new naming convention service.
// Parameters 'context', 'region' and 'environment' are mandatory.
// Context should be set to the application name. In case of a complex system with multiple modules / subsystems,
// 'module' should be set as well.
// Environment provides you the possibility to isolate your resources per environment, e.g. prod, int, stage, dev, etc.
// Regions should be in form of Azure long regions string, e.g. westeurope or northeurope.
// Per default, normal naming convention pattern is used. If required, you can override the pattern using the
// SetPattern() method
func New(context string, region string, environment string, module string) (*Service, error) {
	if context == "" {
		return nil, errors.New("context must be provided")
	}

	if region == "" {
		return nil, errors.New("region must be provided")
	}

	if environment == "" {
		return nil, errors.New("environment must be provided")
	}

	return &Service{
		pattern:     patterns.Normal,
		context:     context,
		module:      module,
		region:      region,
		environment: environment,
	}, nil
}

// SetPattern changes the naming convention pattern to a user defined value. To set a custom pattern, combine the
// placeholders in any order you wish, but without spaces, hyphens or any other characters. Placeholders supported are
// {context} {module} {name} {region} {environment} and {resource_suffix}. For example, you can declare a new pattern
// like this: var myPattern patterns.Pattern = "{resource_suffix}{environment}{context}{module}{region}{name}"
func (service *Service) SetPattern(pattern patterns.Pattern) error {
	numberOfPlaceholders := 6
	mandatoryPlaceholders := []string{"{resource_suffix}", "{environment}", "{context}", "{module}", "{region}", "{name}"}

	for _, placeholder := range mandatoryPlaceholders {
		if !strings.Contains(string(pattern), placeholder) {
			return NewNamingError("placeholder " + placeholder + " not found in the pattern")
		}

		if strings.Count(string(pattern), placeholder) != 1 {
			return NewNamingError("multiple occurrences of " + placeholder + " found in the pattern")
		}
	}

	// to prove that no funny characters are contained in the pattern, we can simply check that the }{ combination occurs
	// the fixed amount of times, which we know because we know how many placeholders are there. Also, we know how the pattern should
	// always start and end.
	if strings.Count(string(pattern), "}{") != (numberOfPlaceholders-1) ||
		!strings.HasPrefix(string(pattern), "{") ||
		!strings.HasSuffix(string(pattern), "}") {
		return NewNamingError("invalid characters found in the pattern, make sure only the placeholders and no other characters are specified")
	}

	service.pattern = pattern
	return nil
}

// GenerateResourceName generates a full Azure resource name, based on the configured naming convention pattern.
// Name parameter should be used to uniquely isolate the resource, in cases where multiple resources on same type
// exist in the same context / module. Name parameter can also be left empty, in which case it will be omitted from the
// pattern during the generation.
func (service *Service) GenerateResourceName(resourceType resources.AzureResourceType, name string) (string, error) {
	var currentNamingConvention = findNamingConvention(resourceType)

	abbreviatedRegion := regions.GetAbbreviatedRegion(service.region)

	namingParts := make([]string, 0)

	type placeholderMapping struct {
		placeholder string
		position    int
		value       string
	}

	placeholderMappings := []placeholderMapping{
		{"{context}", 0, service.context},
		{"{module}", 0, service.module},
		{"{name}", 0, name},
		{"{region}", 0, abbreviatedRegion},
		{"{environment}", 0, service.environment},
		{"{resource_suffix}", 0, string(resourceType)},
	}

	// if name not provided, we need to omit the {name} field from the naming completely
	if name == "" {
		linq.From(placeholderMappings).WhereT(func(mapping placeholderMapping) bool {
			return mapping.placeholder != "{name}"
		}).ToSlice(&placeholderMappings)
	}

	for i := range placeholderMappings {
		// 1. we set the position of each placeholder on the string index of where we find it
		placeholderMappings[i].position = strings.Index(string(service.pattern), placeholderMappings[i].placeholder)
	}

	// 2. so that now we can sort them based on position, and output the values as the naming parts
	linq.From(placeholderMappings).OrderByT(func(mapping placeholderMapping) int {
		return mapping.position
	}).SelectT(func(mapping placeholderMapping) string {
		return mapping.value
	}).ToSlice(&namingParts)

	var separator string

	// hyphen should always be preferred since supported by more Azure resources than underscore
	if currentNamingConvention.Hyphen {
		separator = "-"
	} else if currentNamingConvention.Underscore {
		separator = "_"
	}

	var result = strings.Join(namingParts, separator)

	if currentNamingConvention.Case == UpperCase {
		result = strings.ToUpper(result)
	} else if currentNamingConvention.Case == LowerCase {
		result = strings.ToLower(result)
	}

	valid, err := currentNamingConvention.isValid(result)

	if !valid {
		return "", err
	}

	return result, err
}
