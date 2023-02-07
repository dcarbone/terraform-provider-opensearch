package provider

import (
	"fmt"
	"os"
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
			ProtoV6ProviderFactories: providerFactories,
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
		if os.Getenv("OPENSEARCH_USERNAME") == "" {
			t.Setenv("OPENSEARCH_USERNAME", "admin")
		}
		if os.Getenv("OPENSEARCH_PASSWORD") == "" {
			t.Setenv("OPENSEARCH_PASSWORD", "admin")
		}
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: providerFactories,
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

	t.Run("nested-attributes", func(t *testing.T) {
		const (
			backendRole1 = "backend_role_1"
			backendRole2 = "backend_role_2"

			indexPattern1 = "index_pattern_1*"
			indexPattern2 = "index_*_2"

			dlsValue = "dls_value"
			flsValue = "fls_value"

			maskedField1 = "masked_field_1"
			maskedField2 = "masked_field_2"

			allowedAction1 = "indices:admin/create"
			allowedAction2 = "indices:read*"
		)

		if os.Getenv("OPENSEARCH_USERNAME") == "" {
			t.Setenv("OPENSEARCH_USERNAME", "admin")
		}
		if os.Getenv("OPENSEARCH_PASSWORD") == "" {
			t.Setenv("OPENSEARCH_PASSWORD", "admin")
		}
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: providerFactories,
			Steps: []resource.TestStep{
				{
					Config: acctest.CombineConfig(
						acctest.ProviderConfigLocalhostWith(),
						acctest.PluginSecurityRoleValidConfigWith(
							resourceName,
							map[string]interface{}{
								fields.ResourceAttrBackendRoles: []string{
									backendRole1,
									backendRole2,
								},
								fields.ResourceAttrIndexPermissions: []map[string]interface{}{
									{
										fields.ResourceAttrIndexPatterns: []string{
											indexPattern1,
											indexPattern2,
										},
										fields.ResourceAttrDLS: dlsValue,
										fields.ResourceAttrFLS: flsValue,
										fields.ResourceAttrMaskedFields: []string{
											maskedField1,
											maskedField2,
										},
										fields.ResourceAttrAllowedActions: []string{
											allowedAction1,
											allowedAction2,
										},
									},
								},
							},
						),
					),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceFQN, fields.ResourceAttrRoleName, resourceName),
						resource.TestCheckResourceAttr(
							resourceFQN,
							fmt.Sprintf("%s.0", fields.ResourceAttrBackendRoles),
							backendRole1,
						),
						resource.TestCheckResourceAttr(
							resourceFQN,
							fmt.Sprintf("%s.1", fields.ResourceAttrBackendRoles),
							backendRole2,
						),
					),
				},
			},
		})
	})
}
