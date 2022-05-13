package azure_login

import (
	"github.com/conplementag/cops-hq/pkg/commands"
	"github.com/spf13/viper"
)

type Login struct {
	servicePrincipalId     string
	servicePrincipalSecret string
	tenant                 string
	executor               commands.Executor
}

// New creates a new Login instance by relying on Viper for necessary configuration. Supported
// viper flags are:
// - service-principal-id
// - service-principal-secret
// - service-principal-tenant
func New(executor commands.Executor) *Login {
	return &Login{
		servicePrincipalId:     viper.GetString("service-principal-id"),
		servicePrincipalSecret: viper.GetString("service-principal-secret"),
		tenant:                 viper.GetString("service-principal-tenant"),
		executor:               executor,
	}
}

// NewWithParams creates a new Login instance with the ability to provide all parameters directly
func NewWithParams(executor commands.Executor, servicePrincipalId string, servicePrincipalSecret string, tenant string) *Login {
	return &Login{
		servicePrincipalId:     servicePrincipalId,
		servicePrincipalSecret: servicePrincipalSecret,
		tenant:                 tenant,
		executor:               executor,
	}
}
