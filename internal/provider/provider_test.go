package provider

import (
	"testing"

	"github.com/dcarbone/terraform-provider-opensearch/internal/acctest"
	"github.com/dcarbone/terraform-provider-opensearch/internal/fields"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var providerFactories = map[string]func() (tfprotov6.ProviderServer, error){
	fields.ProviderName: providerserver.NewProtocol6WithError(NewOpenSearchProvider("test")()),
}

func testAccPreCheck(t *testing.T) {
	// You can add code here to run prior to any test case execution, for example assertions
	// about the appropriate environment variables being set are common to see in a pre-check
	// function.
}

func TestBuild_Provider(t *testing.T) {
	_ = NewOpenSearchProvider("test")
}

func TestUnit_ProviderConfig(t *testing.T) {
	t.Run("env", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: providerFactories,
			Steps: []resource.TestStep{
				{
					Config: acctest.ProviderConfigWith(),
				},
			},
		})
	})

	t.Run("static", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: providerFactories,
			Steps: []resource.TestStep{
				{
					Config: acctest.ProviderConfigLocalhostWith(nil),
				},
			},
		})
	})

	t.Run("static-skip-init", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: providerFactories,
			Steps: []resource.TestStep{
				{
					Config: acctest.ProviderConfigLocalhostWith(map[string]interface{}{
						fields.ConfigAttrSkipInitProductCheck: true,
					}),
				},
			},
		})
	})

	t.Run("static-with-logger", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: providerFactories,
			Steps: []resource.TestStep{
				{
					Config: acctest.ProviderConfigLocalhostWith(map[string]interface{}{
						fields.ConfigAttrRequestTraceLogger: map[string]interface{}{
							fields.ConfigAttrEnabled: true,
						},
					}),
				},
			},
		})
	})
}
