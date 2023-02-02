package provider

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/dcarbone/terraform-provider-opensearch/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/opensearch-project/opensearch-go"
	"github.com/opensearch-project/opensearch-go/opensearchapi"
)

func fetchRoles(ctx context.Context, osClient *opensearch.Client, roleName string, diags diag.Diagnostics) (*opensearchapi.Response, client.PluginSecurityRolesAPIResponse) {
	// init opensearch request
	osReq := client.PluginSecurityRolesGetRequest{
		Name: roleName,
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	osResp, err := osReq.Do(ctx, osClient)
	if err != nil {
		diags.AddError(
			"Error querying for role",
			fmt.Sprintf("Error occured querying for role %q: %v", roleName, err),
		)
		return osResp, nil
	}

	// attempt to decode response
	roleResp := make(client.PluginSecurityRolesAPIResponse)
	if err = client.ParseResponse(osResp, &roleResp, http.StatusOK); err != nil {
		if m, ok := err.(client.APIResponseMeta); ok {
			m.AppendDiagnostics(diags)
			return osResp, nil
		}
	}

	return osResp, roleResp
}
