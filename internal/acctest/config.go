package acctest

import (
	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/dcarbone/terraform-provider-opensearch/internal/fields"
)

func ProviderConfigEmpty() string {
	return at.CompileProviderConfig(
		fields.ProviderName,
		map[string]interface{}{},
	)
}

func ProviderConfigLocalhostWith(extra map[string]interface{}) string {
	return at.CompileProviderConfig(
		fields.ProviderName,
		map[string]interface{}{
			fields.ConfigAttrAddresses:             "http://127.0.0.1:9200",
			fields.ConfigAttrInsecureSkipTLSVerify: true,
			fields.ConfigAttrUsername:              "admin",
			fields.ConfigAttrPassword:              "admin",
		},
		extra,
	)
}
