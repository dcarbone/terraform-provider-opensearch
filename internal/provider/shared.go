package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/opensearch-project/opensearch-go"
)

type Shared struct {
	Client *opensearch.Client
}

type ResourceShared struct {
	providerTypeName string
	client           *opensearch.Client
}

func (s *ResourceShared) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		resp.Diagnostics.AddWarning("Provider is not configured", "Provider is not configured")
		return
	}

	// ensure we got what we expected
	shd, ok := req.ProviderData.(*Shared)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected %T, got: %T. Please report this issue to the provider developers.",
				new(Shared),
				req.ProviderData,
			),
		)

		return
	}

	// embed client
	s.client = shd.Client
}
