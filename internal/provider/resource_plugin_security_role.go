package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dcarbone/terraform-provider-opensearch/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/opensearch-project/opensearch-go/opensearchapi"
)

func NewPluginSecurityRoleResource() resource.Resource {
	r := new(PluginSecurityRoleResource)
	return r
}

type PluginSecurityRoleResource struct {
	ResourceShared
}

type PluginSecurityRoleResourceData struct {
	RoleName           types.String `tfsdk:"role_name"`
	Description        types.String `tfsdk:"description"`
	ClusterPermissions types.List   `tfsdk:"cluster_permissions"`
	IndexPermissions   types.List   `tfsdk:"index_permissions"`
	TenantPermissions  types.List   `tfsdk:"tenant_permissions"`

	Reserved types.Bool `tfsdk:"reserved"`
	Hidden   types.Bool `tfsdk:"hidden"`
	Static   types.Bool `tfsdk:"static"`
}

func (r *PluginSecurityRoleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = makeResourceName(req.ProviderTypeName, resourceSuffixSecurityPluginRole)
}

func (r *PluginSecurityRoleResource) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "OpenSearch Security Plugin Role",
		Attributes: map[string]schema.Attribute{
			resourceAttrRoleName: schema.StringAttribute{
				Required: true,
			},
			resourceAttrDescription: schema.StringAttribute{
				Optional: true,
			},
			resourceAttrClusterPermissions: schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
			resourceAttrIndexPermissions: schema.ListNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						resourceAttrIndexPatterns: schema.ListAttribute{
							Optional:    true,
							ElementType: types.StringType,
						},
						resourceAttrDLS: schema.StringAttribute{
							Optional: true,
						},
						resourceAttrFLS: schema.StringAttribute{
							Optional: true,
						},
						resourceAttrMaskedFields: schema.ListAttribute{
							Optional:    true,
							ElementType: types.StringType,
						},
						resourceAttrAllowedActions: schema.ListAttribute{
							Optional:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
			resourceAttrTenantPermissions: schema.ListNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						resourceAttrTenantPatterns: schema.ListAttribute{
							Optional:    true,
							ElementType: types.StringType,
						},
						resourceAttrAllowedActions: schema.ListAttribute{
							Optional:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
			resourceAttrStatic: schema.BoolAttribute{
				Computed: true,
			},
			resourceAttrHidden: schema.BoolAttribute{
				Computed: true,
			},
			resourceAttrReserved: schema.BoolAttribute{
				Computed: true,
			},
		},
	}
}

func (r *PluginSecurityRoleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	role := new(PluginSecurityRoleResourceData)

	resp.Diagnostics.Append(req.Plan.Get(ctx, role)...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (r *PluginSecurityRoleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var (
		roleName string
		osReq    *client.PluginSecurityRoleGetRequest
		osResp   *opensearchapi.Response
		osRole   client.PluginSecurityRole
		osRoles  client.PluginSecurityRoleList
		ok       bool
		err      error

		stateData = new(PluginSecurityRoleResourceData)
	)

	// marshal state into data type, appending errors to response
	resp.Diagnostics.Append(req.State.Get(ctx, stateData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// extract role name
	roleName = stateData.RoleName.ValueString()

	// prepare opensearch api request
	osReq = &client.PluginSecurityRoleGetRequest{
		Name: roleName,
	}

	// construct request context
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// execute query, and set up response body cleanup
	osResp, err = osReq.Do(ctx, r.client)
	defer client.HandleResponseCleanup(osResp)

	// check for client error
	if err != nil {
		resp.Diagnostics.AddError(
			"Error executing OpenSearch API query",
			fmt.Sprintf("Received error querying for role %q: %v", roleName, err),
		)
		return
	}

	// decode response
	if err = json.NewDecoder(osResp.Body).Decode(&osRoles); err != nil {
		resp.Diagnostics.AddError(
			"Error decoding OpenSearch Role response",
			fmt.Sprintf("Error decoding OpenSearch Role: %v", err),
		)
		return
	}

	// check for api errors
	if osRoles.HasErrors() {
		osRoles.AddDiagnosticErrors(resp.Diagnostics)
		return
	}

	// locate specific role
	if osRole, ok = osRoles.Roles[roleName]; !ok {
		resp.Diagnostics.AddError(
			"Role not found",
			fmt.Sprintf("Role %q not found", roleName),
		)
		return
	}

	if osRole.Hidden != nil {
		stateData.Hidden = types.BoolValue(*osRole.Hidden)
	} else {
		stateData.Hidden = types.BoolNull()
	}
}

func (r *PluginSecurityRoleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

}

func (r *PluginSecurityRoleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

}

func (r *PluginSecurityRoleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

}
