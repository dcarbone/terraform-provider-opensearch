package provider

import (
	"context"
	"net/http"

	"github.com/dcarbone/terraform-provider-opensearch/internal/client"
	"github.com/opensearch-project/opensearch-go"
	"github.com/opensearch-project/opensearch-go/opensearchapi"
)

func tryFetchRoles(ctx context.Context, osClient *opensearch.Client, roleName string) (*opensearchapi.Response, client.PluginSecurityRolesAPIResponse, error) {
	// init opensearch request
	osReq := client.PluginSecurityRolesGetRequest{
		Name: roleName,
	}

	osResp, err := osReq.Do(ctx, osClient)
	if err != nil {
		// just in case response isn't nil
		defer client.HandleResponseCleanup(osResp)
		return nil, nil, err
	}

	// attempt to decode response
	roleResp := make(client.PluginSecurityRolesAPIResponse)
	if err = client.ParseResponse(osResp, &roleResp, http.StatusOK); err != nil {
		return osResp, nil, err
	}

	return osResp, roleResp, nil
}
