package naming

import (
	"github.com/ahmetb/go-linq"
	"github.com/conplementag/cops-hq/pkg/naming/resources"
	"regexp"
	"strconv"
	"strings"
)

type namingConvention struct {
	Type         resources.AzureResourceType
	MinLength    int
	MaxLength    int
	Alphanumeric bool
	Hyphen       bool
	Underscore   bool
	Case         caseSensitivity
}

var namingConventions = []namingConvention{
	{resources.ResourceGroup, 3, 90, true, true, true, CaseInsensitive},
	{resources.SqlServer, 3, 63, true, true, false, LowerCase},
	{resources.SqlDatabase, 3, 63, true, true, false, LowerCase},
	{resources.SqlManagedInstance, 3, 63, true, true, false, LowerCase},
	{resources.SqlElasticPool, 3, 63, true, true, false, LowerCase},
	{resources.KeyVault, 3, 24, true, true, false, CaseInsensitive},
	{resources.IotHub, 3, 50, true, true, false, CaseInsensitive},
	{resources.RecoveryServicesVault, 5, 50, true, true, false, CaseInsensitive},
	{resources.AKSCluster, 5, 50, true, true, true, CaseInsensitive},
	{resources.StorageAccount, 3, 24, true, false, false, LowerCase},
	{resources.VirtualNetwork, 3, 64, true, true, true, CaseInsensitive},
	{resources.VirtualNetworkGateway, 3, 64, true, true, true, CaseInsensitive},
	{resources.RouteTable, 3, 80, true, true, true, CaseInsensitive},
	{resources.ApplicationGateway, 3, 90, true, true, true, CaseInsensitive},
	{resources.PublicIp, 3, 90, true, true, true, CaseInsensitive},
	{resources.Bastion, 3, 64, true, true, true, CaseInsensitive},
	{resources.UserAssignedIdentity, 3, 90, true, true, true, CaseInsensitive},
	{resources.NetworkSecurityGroup, 3, 80, true, true, true, CaseInsensitive},
	{resources.LogAnalyticsWorkspace, 3, 90, false, true, true, CaseInsensitive},
}

func findNamingConvention(resourceType resources.AzureResourceType) namingConvention {
	return linq.From(namingConventions).SingleWithT(func(c namingConvention) bool {
		return c.Type == resourceType
	}).(namingConvention)
}

// isValid verifies the naming convention against a given value parameter
func (r namingConvention) isValid(value string) (bool, error) {
	if len(value) < r.MinLength {
		return false, NewNamingError("Min length not reached")
	}

	if len(value) > r.MaxLength {
		return false, NewNamingError("Max length exceeded")
	}

	if regexp.MustCompile(r.getRegexPattern(len(value))).MatchString(value) == false {
		return false, NewNamingError("Invalid char used")
	}

	if r.Case == UpperCase && strings.ToUpper(value) != value {
		return false, NewNamingError("Must not contain lowercase characters")
	}

	if r.Case == LowerCase && strings.ToLower(value) != value {
		return false, NewNamingError("Must not contain uppercase characters")
	}

	return true, nil
}

func (r namingConvention) getRegexPattern(length int) string {
	var patterns []string
	patterns = append(patterns, "\\w")
	if r.Underscore {
		patterns = append(patterns, "_")
	}
	if r.Hyphen {
		patterns = append(patterns, "-")
	}

	return "[" + strings.Join(patterns, ",") + "]{" + strconv.Itoa(length) + "," + strconv.Itoa(length) + "}?"
}
