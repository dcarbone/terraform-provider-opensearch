package provider

import (
	"regexp"
	"testing"

	"github.com/dcarbone/terraform-provider-opensearch/internal/acctest"
	"github.com/dcarbone/terraform-provider-opensearch/internal/fields"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_PluginSecurityRole(t *testing.T) {
	const (
		resourceName = "test_role"
	)

	var (
		resourceFQN = fields.ResourceTypeFQN(fields.ProviderName, fields.ResourceTypeSecurityPluginRole, resourceName)
	)

	t.Run("empty-throws-error", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: protoV6ProviderFactories,
			Steps: []resource.TestStep{ // Create and Read testing
				{
					Config: acctest.CombineConfig(
						acctest.ProviderConfigLocalhostWith(),
						acctest.PluginSecurityRoleConfigWith(resourceName, nil),
					),
					ExpectError: regexp.MustCompile("required"),
				},
			},
		})
	})

	t.Run("basic", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: protoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: acctest.CombineConfig(
						acctest.ProviderConfigLocalhostWith(),
						acctest.PluginSecurityRoleValidConfigWith(resourceName, nil),
					),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceFQN, fields.ResourceAttrRoleName, resourceName),
					),
				},
			},
		})
	})

	//resource.Test(t, resource.TestCase{
	//	PreCheck:                 func() { testAccPreCheck(t) },
	//	ProtoV6ProviderFactories: protoV6ProviderFactories,
	//	Steps: []resource.TestStep{
	//		// Create and Read testing
	//		{
	//			Config: acctest.PluginSecurityRoleConfigWith(resourceName, nil),
	//			Check: resource.ComposeAggregateTestCheckFunc(
	//				resource.TestCheckResourceAttr("scaffolding_example.test", "configurable_attribute", "one"),
	//				resource.TestCheckResourceAttr("scaffolding_example.test", "id", "example-id"),
	//			),
	//		},
	//		// ImportState testing
	//		{
	//			ResourceName:      "scaffolding_example.test",
	//			ImportState:       true,
	//			ImportStateVerify: true,
	//			// This is not normally necessary, but is here because this
	//			// example code does not have an actual upstream service.
	//			// Once the Read method is able to refresh information from
	//			// the upstream service, this can be removed.
	//			ImportStateVerifyIgnore: []string{"configurable_attribute"},
	//		},
	//		// Update and Read testing
	//		{
	//			Config: testAccExampleResourceConfig("two"),
	//			Check: resource.ComposeAggregateTestCheckFunc(
	//				resource.TestCheckResourceAttr("scaffolding_example.test", "configurable_attribute", "two"),
	//			),
	//		},
	//		// Delete testing automatically occurs in TestCase
	//	},
	//})
}
