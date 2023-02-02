package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dcarbone/terraform-plugin-framework-utils/v3/conv"
	"github.com/dcarbone/terraform-provider-opensearch/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/diag"
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

func (d *PluginSecurityRoleResourceData) UpdateFromRole(roleName string, r client.PluginSecurityRole) diag.Diagnostics {
	var diags diag.Diagnostics

	d.RoleName = types.StringValue(roleName)
	d.Description = types.StringValue(r.Description)
	d.ClusterPermissions = conv.StringsToStringList(r.ClusterPermissions, true)

	if d.IndexPermissions, diags = indexPermissionsToTerraformNestedList(r.IndexPermissions, false); diags.HasError() {
		return diags
	}
	if d.TenantPermissions, diags = tenantPermissionsToTerraformNestedList(r.TenantPermissions, false); diags.HasError() {
		return diags
	}

	// set "computed" values
	d.Hidden = conv.BoolPtrToBoolValue(r.Hidden)
	d.Static = conv.BoolPtrToBoolValue(r.Static)
	d.Reserved = conv.BoolPtrToBoolValue(r.Reserved)

	return diag.Diagnostics{}
}

func (r *PluginSecurityRoleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = makeResourceName(req.ProviderTypeName, resourceSuffixSecurityPluginRole)
}

func (r *PluginSecurityRoleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
	var (
		roleName string
		osRole   client.PluginSecurityRole
		osReq    *client.PluginSecurityRoleUpsertRequest
		osResp   *opensearchapi.Response
		jsonB    []byte
		err      error

		planData = new(PluginSecurityRoleResourceData)
	)

	// marshal plan value into data type, appending errors to response
	resp.Diagnostics.Append(req.Plan.Get(ctx, planData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// extract role name
	roleName = planData.RoleName.ValueString()

	{
		// attempt to locate existing role by name
		osReq := &client.PluginSecurityRoleGetRequest{Name: roleName}
		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		osResp, err := osReq.Do(ctx, r.client)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error executing OpenSearch API query",
				fmt.Sprintf("Received error querying for role %q: %v", roleName, err),
			)
			return
		}

		// if response is anything other than 404, assume big problem
		if osResp.StatusCode != http.StatusNotFound {
			resp.Diagnostics.AddError(
				"Security Role already exists",
				fmt.Sprintf("Cannot create new Security Role %q as it already exists in the cluster or your credentials do not have sufficient permissions", roleName),
			)
			return
		}
	}

	// init request type
	osReq = &client.PluginSecurityRoleUpsertRequest{
		Name: roleName,
	}

	// convert plan data to opensearch model
	osRole = terraformSecurityRoleToSecurityRole(planData)

	if jsonB, err = json.Marshal(osRole); err != nil {
		resp.Diagnostics.AddError(
			"Error marshalling plan into OpenSearch request",
			fmt.Sprintf("Error json-encoding plan data into OpenSearch request: %v", err),
		)
		return
	}

	// set request body
	osReq.Body = bytes.NewReader(jsonB)

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if osResp, err = osReq.Do(ctx, r.client); err != nil {
		resp.Diagnostics.AddError(
			"Error creating role",
			fmt.Sprintf("Error executing create role request: %v", err),
		)
		return
	}

	//planData.UpdateFromRole(roleName, )
}

func (r *PluginSecurityRoleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var (
		roleName    string
		osReq       *client.PluginSecurityRoleGetRequest
		osResp      *opensearchapi.Response
		osRole      client.PluginSecurityRole
		osRoles     client.PluginSecurityRoleList
		updateDiags diag.Diagnostics
		ok          bool
		err         error

		stateData = new(PluginSecurityRoleResourceData)
	)

	// marshal state value into data type, appending errors to response
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
		osRoles.AppendErrorsToDiagnostic(resp.Diagnostics)
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

	// update data object from source role
	updateDiags = stateData.UpdateFromRole(roleName, osRole)

	// append any / all diagnostics to response diags
	resp.Diagnostics.Append(updateDiags...)

	// if there were errors, end now
	if resp.Diagnostics.HasError() {
		return
	}

	// update state from remote, appending any resulting diags
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *PluginSecurityRoleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

}

func (r *PluginSecurityRoleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

}

func (r *PluginSecurityRoleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

}
