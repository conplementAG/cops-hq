package naming

import (
	"github.com/conplementag/cops-hq/pkg/naming/patterns"
	"github.com/conplementag/cops-hq/pkg/naming/resources"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_GenerateResourceName(t *testing.T) {
	type args struct {
		name         string
		pattern      patterns.Pattern
		resourceType resources.AzureResourceType
	}

	namingService, err := New("acme", "westeurope", "dev", "front")
	assert.NoError(t, err)

	var newPattern patterns.Pattern = "{resource_suffix}{module}{name}{environment}{region}{context}"

	tests := []struct {
		testName       string
		expectedResult string
		expectedError  error
		args           args
	}{
		{"normal azure resource", "acme-front-green-weu-dev-rg", nil,
			args{"green", patterns.Normal, resources.ResourceGroup},
		},
		{"short length azure resource", "acmefrontgreenweudevsa", nil,
			args{"green", patterns.Normal, resources.StorageAccount},
		},
		{"short length azure resource - too long", "", NewNamingError("Max length exceeded"),
			args{"alongname", patterns.Normal, resources.StorageAccount},
		},
		{"invalid char used", "", NewNamingError("Invalid char used"),
			args{"la&la", patterns.Normal, resources.StorageAccount},
		},
		{"lowercase test", "acme-front-thename-weu-dev-sqls", nil,
			args{"TheName", patterns.Normal, resources.SqlServer},
		},
		{"changed pattern", "sqls-front-green-dev-weu-acme", nil,
			args{"green", newPattern, resources.SqlServer},
		},
		{"name can be omitted", "acmefrontweudevsa", nil,
			args{"", patterns.Normal, resources.StorageAccount},
		},
		{"changed pattern with omitted name", "sqls-front-dev-weu-acme", nil,
			args{"", newPattern, resources.SqlServer},
		},
	}

	for _, tt := range tests {
		namingService.SetPattern(tt.args.pattern)
		got, err := namingService.GenerateResourceName(tt.args.resourceType, tt.args.name)

		if got != tt.expectedResult {
			t.Errorf("GenerateResourceName() got = %v, expected %v", got, tt.expectedResult)
		}

		if tt.expectedError == nil {
			assert.NoError(t, err)
		} else {
			if assert.Error(t, err) {
				assert.Equal(t, tt.expectedError, err)
			}
		}
	}
}

func Test_SetPattern(t *testing.T) {
	namingService, err := New("acme", "westeurope", "dev", "front")
	assert.NoError(t, err)

	var customPattern patterns.Pattern = "{resource_suffix}{module}{name}{environment}{region}{context}"
	var duplicatesPattern patterns.Pattern = "{resource_suffix}{module}{name}{environment}{region}{context}{name}"
	var missingPartsPattern patterns.Pattern = "{resource_suffix}{module}{name}{region}{context}"
	var invalidCharsPattern patterns.Pattern = "{resource_suffix}{module}-{name}{environment}_{region}{context}"

	tests := []struct {
		testName      string
		pattern       patterns.Pattern
		expectedError error
	}{
		{"normal pattern", patterns.Normal, nil},
		{"custom pattern", customPattern, nil},
		{"pattern with duplicates", duplicatesPattern, NewNamingError("multiple occurrences of {name} found in the pattern")},
		{"pattern with missing parts", missingPartsPattern, NewNamingError("placeholder {environment} not found in the pattern")},
		{"pattern with invalid chars", invalidCharsPattern,
			NewNamingError("invalid characters found in the pattern, make sure only the placeholders and no other characters are specified")},
	}

	for _, tt := range tests {
		err := namingService.SetPattern(tt.pattern)

		if tt.expectedError == nil {
			assert.NoError(t, err)
		} else {
			if assert.Error(t, err) {
				assert.Equal(t, tt.expectedError, err)
			}
		}
	}
}
