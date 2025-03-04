package azure_login

import (
	"github.com/conplementag/cops-hq/v2/pkg/commands"
	"github.com/spf13/viper"
)

// New creates a new Login instance by relying on Viper for necessary configuration. Supported
// viper flags are:
//   - service-principal-id
//   - service-principal-secret
//   - service-principal-tenant
//   - user-assigned-managed-identity-client-id
//   - managed-identity-tenant-id
//   - use-managed-identity
func New(executor commands.Executor) *Login {
	return &Login{
		servicePrincipalId:                  viper.GetString("service-principal-id"),
		servicePrincipalSecret:              viper.GetString("service-principal-secret"),
		servicePrincipalTenantId:            viper.GetString("service-principal-tenant"),
		userAssignedManagedIdentityClientId: viper.GetString("user-assigned-managed-identity-client-id"),
		managedIdentityTenantId:             viper.GetString("managed-identity-tenant-id"),
		useManagedIdentity:                  viper.GetBool("use-managed-identity"),
		executor:                            executor,
	}
}

// NewWithParams creates a new Login instance with the ability to provide all parameters directly
func NewWithParams(executor commands.Executor, servicePrincipalId string, servicePrincipalSecret string, servicePrincipalTenantId string, userAssignedManagedIdentityClientId string, managedIdentityTenantId string, useManagedIdentity bool) *Login {
	return &Login{
		servicePrincipalId:                  servicePrincipalId,
		servicePrincipalSecret:              servicePrincipalSecret,
		servicePrincipalTenantId:            servicePrincipalTenantId,
		userAssignedManagedIdentityClientId: userAssignedManagedIdentityClientId,
		managedIdentityTenantId:             managedIdentityTenantId,
		useManagedIdentity:                  useManagedIdentity,
		executor:                            executor,
	}
}
