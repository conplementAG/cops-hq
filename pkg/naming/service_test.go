package naming

import (
	"github.com/conplementag/cops-hq/v2/pkg/naming/patterns"
	"github.com/conplementag/cops-hq/v2/pkg/naming/resources"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_GenerateResourceName(t *testing.T) {
	type args struct {
		module       string
		name         string
		color        string
		pattern      patterns.Pattern
		resourceType resources.AzureResourceType
	}

	var newPattern patterns.Pattern = "{resource_suffix}{module}{color}{name}{environment}{region}{context}"

	tests := []struct {
		testName       string
		expectedResult string
		expectedError  error
		args           args
	}{
		{"normal azure resource", "acme-front-g-bla-weu-dev-rg", nil,
			args{"front", "bla", "g", patterns.Normal, resources.ResourceGroup},
		},
		{"short length azure resource", "acmefrontgblaweudevsa", nil,
			args{"front", "bla", "g", patterns.Normal, resources.StorageAccount},
		},
		{"short length azure resource - too long", "", NewNamingError("Max length exceeded"),
			args{"front", "alongname", "b", patterns.Normal, resources.StorageAccount},
		},
		{"invalid char used", "", NewNamingError("Invalid char used"),
			args{"front", "la&la", "b", patterns.Normal, resources.StorageAccount},
		},
		{"lowercase test", "acme-front-b-thename-weu-dev-sqls", nil,
			args{"front", "TheName", "b", patterns.Normal, resources.SqlServer},
		},
		{"changed pattern", "sqls-front-b-bla-dev-weu-acme", nil,
			args{"front", "bla", "b", newPattern, resources.SqlServer},
		},
		{"name can be omitted", "acmefrontgweudevsa", nil,
			args{"front", "", "g", patterns.Normal, resources.StorageAccount},
		},
		{"name can be omitted when resource name with dashes requested", "acme-front-g-weu-dev-rg", nil,
			args{"front", "", "g", patterns.Normal, resources.ResourceGroup},
		},
		{"changed pattern with omitted name", "sqls-front-g-dev-weu-acme", nil,
			args{"front", "", "g", newPattern, resources.SqlServer},
		},
		{"module can be omitted - normal azure resource", "acme-g-test-weu-dev-rg", nil,
			args{"", "test", "g", patterns.Normal, resources.ResourceGroup},
		},
		{"module can be omitted - short length azure resource", "acmextestweudevsa", nil,
			args{"", "test", "x", patterns.Normal, resources.StorageAccount},
		},
		{"long resource name test storage account", "acmefrontbtestweudevsa", nil,
			args{"front", "test", "b", patterns.Normal, resources.StorageAccount},
		},
		{"long resource name test private endpoint", "acme-front-b-test-weu-dev-pe", nil,
			args{"front", "test", "b", patterns.Normal, resources.PrivateEndpoint},
		},
		{"long resource name test key vault", "acmefrontbtestweudevkv", nil,
			args{"front", "test", "b", patterns.Normal, resources.KeyVault},
		},
		{"long resource name test key vault", "acme-front-b-test-weu-dev-sqlmi", nil,
			args{"front", "test", "b", patterns.Normal, resources.SqlManagedInstance},
		},
	}

	for _, tt := range tests {
		namingService, err := New("acme", tt.args.module, tt.args.color, "westeurope", "dev")
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
	namingService, err := New("acme", "front", "blue", "westeurope", "dev")
	assert.NoError(t, err)

	var customPattern patterns.Pattern = "{resource_suffix}{module}{color}{name}{environment}{region}{context}"
	var duplicatesPattern patterns.Pattern = "{resource_suffix}{module}{name}{color}{environment}{region}{context}{name}"
	var missingPartsPattern patterns.Pattern = "{resource_suffix}{module}{name}{region}{context}"
	var invalidCharsPattern patterns.Pattern = "{resource_suffix}{module}-{name}{color}{environment}_{region}{context}"

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
