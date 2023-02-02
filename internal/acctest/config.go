package acctest

import (
	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/dcarbone/terraform-provider-opensearch/internal/fields"
	"strings"
)

func CombineConfig(in ...string) string {
	return strings.Join(in, "\n\n")
}

func ProviderConfigWith(extra ...map[string]interface{}) string {
	return at.CompileProviderConfig(
		fields.ProviderName,
		extra...,
	)
}

func ProviderConfigLocalhostWith(extra ...map[string]interface{}) string {
	return ProviderConfigWith(
		append(
			[]map[string]interface{}{
				{

					fields.ConfigAttrAddresses:             []string{"https://127.0.0.1:9200"},
					fields.ConfigAttrInsecureSkipTLSVerify: true,
					fields.ConfigAttrUsername:              "admin",
					fields.ConfigAttrPassword:              "admin",
					fields.ConfigAttrLogging: map[string]interface{}{
						fields.ConfigAttrEnabled:             true,
						fields.ConfigAttrIncludeRequestBody:  true,
						fields.ConfigAttrIncludeResponseBody: true,
					},
				},
			},
			extra...,
		)...,
	)
}

func PluginSecurityRoleConfigWith(name string, extra ...map[string]interface{}) string {
	return at.CompileResourceConfig(
		fields.TypeName(fields.ProviderName, fields.ResourceTypeSecurityPluginRole),
		name,
		extra...,
	)
}

func PluginSecurityRoleValidConfigWith(name string, extra ...map[string]interface{}) string {
	return PluginSecurityRoleConfigWith(
		name,
		append(
			[]map[string]interface{}{
				{
					fields.ResourceAttrRoleName: name,
				},
			},
			extra...,
		)...,
	)
}
