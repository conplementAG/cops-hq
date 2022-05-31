package naming

import (
	"github.com/conplementag/cops-hq/pkg/naming/patterns"
	"github.com/conplementag/cops-hq/pkg/naming/resources"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_GenerateResourceName(t *testing.T) {
	type args struct {
		module       string
		name         string
		pattern      patterns.Pattern
		resourceType resources.AzureResourceType
	}

	var newPattern patterns.Pattern = "{resource_suffix}{module}{name}{environment}{region}{context}"

	tests := []struct {
		testName       string
		expectedResult string
		expectedError  error
		args           args
	}{
		{"normal azure resource", "acme-front-green-weu-dev-rg", nil,
			args{"front", "green", patterns.Normal, resources.ResourceGroup},
		},
		{"short length azure resource", "acmefrontgreenweudevsa", nil,
			args{"front", "green", patterns.Normal, resources.StorageAccount},
		},
		{"short length azure resource - too long", "", NewNamingError("Max length exceeded"),
			args{"front", "alongname", patterns.Normal, resources.StorageAccount},
		},
		{"invalid char used", "", NewNamingError("Invalid char used"),
			args{"front", "la&la", patterns.Normal, resources.StorageAccount},
		},
		{"lowercase test", "acme-front-thename-weu-dev-sqls", nil,
			args{"front", "TheName", patterns.Normal, resources.SqlServer},
		},
		{"changed pattern", "sqls-front-green-dev-weu-acme", nil,
			args{"front", "green", newPattern, resources.SqlServer},
		},
		{"name can be omitted", "acmefrontweudevsa", nil,
			args{"front", "", patterns.Normal, resources.StorageAccount},
		},
		{"name can be omitted when resource name with dashes requested", "acme-front-weu-dev-rg", nil,
			args{"front", "", patterns.Normal, resources.ResourceGroup},
		},
		{"changed pattern with omitted name", "sqls-front-dev-weu-acme", nil,
			args{"front", "", newPattern, resources.SqlServer},
		},
		{"module can be omitted - normal azure resource", "acme-green-weu-dev-rg", nil,
			args{"", "green", patterns.Normal, resources.ResourceGroup},
		},
		{"module can be omitted - short length azure resource", "acmegreenweudevsa", nil,
			args{"", "green", patterns.Normal, resources.StorageAccount},
		},
	}

	for _, tt := range tests {
		namingService, err := New("acme", "westeurope", "dev", tt.args.module)
		assert.NoError(t, err)

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
