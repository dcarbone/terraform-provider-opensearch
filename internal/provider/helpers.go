package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dcarbone/terraform-provider-opensearch/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/opensearch-project/opensearch-go"
	"github.com/opensearch-project/opensearch-go/opensearchapi"
)

func fetchRoles(ctx context.Context, osClient *opensearch.Client, roleName string, diags diag.Diagnostics) (*client.PluginSecurityRolesAPIResponse, *opensearchapi.Response) {
	// init opensearch request
	osReq := client.PluginSecurityRolesGetRequest{
		Name: roleName,
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	osResp, err := osReq.Do(ctx, osClient)

	defer client.HandleResponseCleanup(osResp)

	if err != nil {
		diags.AddError(
			"Error querying for role",
			fmt.Sprintf("Error occured querying for role %q: %v", roleName, err),
		)
		return nil, osResp
	}

	roleResp := new(client.PluginSecurityRolesAPIResponse)
	if err := json.NewDecoder(osResp.Body).Decode(roleResp); err != nil {
		diags.AddError(
			"Error unmarshalling role",
			fmt.Sprintf("Error occurred while unmarshalling API response: %v", err),
		)
		return nil, osResp
	}

	return roleResp, osResp
}
