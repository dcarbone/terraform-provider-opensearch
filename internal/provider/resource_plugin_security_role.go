package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dcarbone/terraform-plugin-framework-utils/v3/conv"
	"github.com/dcarbone/terraform-plugin-framework-utils/v3/validation"
	"github.com/dcarbone/terraform-provider-opensearch/internal/client"
	"github.com/dcarbone/terraform-provider-opensearch/internal/fields"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
	resp.TypeName = fields.TypeName(req.ProviderTypeName, fields.ResourceTypeSecurityPluginRole)
}

func (r *PluginSecurityRoleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "OpenSearch Security Plugin Role",
		Attributes: map[string]schema.Attribute{
			fields.ResourceAttrRoleName: schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					validation.Required(),
				},
			},
			fields.ResourceAttrDescription: schema.StringAttribute{
				Optional: true,
			},
			fields.ResourceAttrClusterPermissions: schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
			fields.ResourceAttrIndexPermissions: schema.ListNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						fields.ResourceAttrIndexPatterns: schema.ListAttribute{
							Optional:    true,
							ElementType: types.StringType,
						},
						fields.ResourceAttrDLS: schema.StringAttribute{
							Optional: true,
						},
						fields.ResourceAttrFLS: schema.StringAttribute{
							Optional: true,
						},
						fields.ResourceAttrMaskedFields: schema.ListAttribute{
							Optional:    true,
							ElementType: types.StringType,
						},
						fields.ResourceAttrAllowedActions: schema.ListAttribute{
							Optional:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
			fields.ResourceAttrTenantPermissions: schema.ListNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						fields.ResourceAttrTenantPatterns: schema.ListAttribute{
							Optional:    true,
							ElementType: types.StringType,
						},
						fields.ResourceAttrAllowedActions: schema.ListAttribute{
							Optional:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
			fields.ResourceAttrStatic: schema.BoolAttribute{
				Computed: true,
			},
			fields.ResourceAttrHidden: schema.BoolAttribute{
				Computed: true,
			},
			fields.ResourceAttrReserved: schema.BoolAttribute{
				Computed: true,
			},
		},
	}
}

func (r *PluginSecurityRoleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var (
		roleName string
		osRole   client.PluginSecurityRole
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
		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		psResp, _, _ := tryFetchRoles(ctx, r.client, roleName)

		if psResp != nil {
			// if we got some kind of response from opensearch, test status code
			if psResp.StatusCode == 200 {
				resp.Diagnostics.AddError(
					"Role already exists",
					fmt.Sprintf("Role %q already exists in cluster", roleName),
				)
				return
			}
			// if we get here, assume that the role either does not already exists, or some kind of permission
			// error occurred at this point, allow create attempt to happen.
		} else if err != nil {
			// if an error was seen, assume big badness
			if m, ok := err.(client.APIStatusResponse); ok {
				m.AppendDiagnostics(resp.Diagnostics)
			} else {
				resp.Diagnostics.AddError(
					"Error querying for role",
					fmt.Sprintf("Error occurred looking for existing role %q: %v", roleName, err.Error()),
				)
			}
			return
		}
	}

	// execute create request
	{
		// init request type
		osReq := &client.PluginSecurityRoleUpsertRequest{
			Name: roleName,
		}

		// convert plan data to opensearch model
		osRole = terraformSecurityRoleToSecurityRole(planData)

		jsonB, err := json.Marshal(osRole)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error marshalling plan into OpenSearch request",
				fmt.Sprintf("Error json-encoding plan data into OpenSearch request: %v", err),
			)
			return
		}

		// set request body
		osReq.Body = bytes.NewReader(jsonB)

		// execute create call
		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		osResp, err := osReq.Do(ctx, r.client)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error creating role",
				fmt.Sprintf("Error executing create role request: %v", err),
			)
			return
		}

		// attempt to parse response
		if err = client.ParseResponse(osResp, &osRole, http.StatusOK); err != nil {
			if m, ok := err.(client.APIStatusResponse); ok {
				m.AppendDiagnostics(resp.Diagnostics)
			} else {
				resp.Diagnostics.AddError(
					"Error parsing create role response",
					err.Error(),
				)
			}
			return
		}
	}

	// otherwise, try to update state model with new data
	resp.Diagnostics.Append(planData.UpdateFromRole(roleName, osRole)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// finally, try to update state itself with updated model
	resp.Diagnostics.Append(resp.State.Set(ctx, planData)...)
}

func (r *PluginSecurityRoleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var (
		roleName    string
		osRole      client.PluginSecurityRole
		updateDiags diag.Diagnostics
		ok          bool

		stateData = new(PluginSecurityRoleResourceData)
	)

	// marshal state value into data type, appending errors to response
	resp.Diagnostics.Append(req.State.Get(ctx, stateData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// extract role name
	roleName = stateData.RoleName.ValueString()

	// query for role from cluster
	// done in sub-context to avoid poisoning ctx var
	{
		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		_, osRoles, err := tryFetchRoles(ctx, r.client, roleName)
		if err != nil {
			if m, ok := err.(client.APIStatusResponse); ok {
				m.AppendDiagnostics(resp.Diagnostics)
			} else {
				resp.Diagnostics.AddError(
					"Error querying for role",
					fmt.Sprintf("Error occurred querying for role %q: %v", roleName, err.Error()),
				)
			}
			return
		}

		// attempt to extract role from response
		if osRole, ok = osRoles[roleName]; !ok {
			resp.Diagnostics.AddError(
				"Role not found",
				fmt.Sprintf("Role %q not found", roleName),
			)
			return
		}
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
	var (
		roleName string
		osRole   client.PluginSecurityRole

		planData = new(PluginSecurityRoleResourceData)
	)

	// marshal plan value into data type, appending errors to response
	resp.Diagnostics.Append(req.Plan.Get(ctx, planData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// extract role name
	roleName = planData.RoleName.ValueString()

	// attempt to locate role in cluster
	{
		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		psResp, _, err := tryFetchRoles(ctx, r.client, roleName)
		if err != nil {
			if m, ok := err.(client.APIStatusResponse); ok {
				m.AppendDiagnostics(resp.Diagnostics)
			} else {
				resp.Diagnostics.AddError(
					"Error querying for role",
					fmt.Sprintf("Error occurred querying for role %q: %v", roleName, err.Error()),
				)
			}
			return
		}

		// if the role was not found, prevent the update call from creating a new one.
		if psResp.StatusCode != 200 {
			resp.Diagnostics.AddError(
				"Role not found",
				fmt.Sprintf("Role %q was not found in cluster", roleName),
			)
			return
		}
	}

	// execute update call
	{
		// init request type
		osReq := &client.PluginSecurityRoleUpsertRequest{
			Name: roleName,
		}

		// convert plan data to opensearch model
		osRole = terraformSecurityRoleToSecurityRole(planData)

		jsonB, err := json.Marshal(osRole)
		if err != nil {
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
		osResp, err := osReq.Do(ctx, r.client)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error creating role",
				fmt.Sprintf("Error executing create role request: %v", err),
			)
			return
		}

		// attempt to parse response
		if err = client.ParseResponse(osResp, &osRole, http.StatusOK); err != nil {
			if m, ok := err.(client.APIStatusResponse); ok {
				m.AppendDiagnostics(resp.Diagnostics)
			} else {
				resp.Diagnostics.AddError(
					"Error parsing create role response",
					err.Error(),
				)
			}
			return
		}
	}

	// otherwise, try to update state model with new data
	resp.Diagnostics.Append(planData.UpdateFromRole(roleName, osRole)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// finally, try to update state itself with updated model
	resp.Diagnostics.Append(resp.State.Set(ctx, planData)...)
}

func (r *PluginSecurityRoleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var (
		roleName string

		planData = new(PluginSecurityRoleResourceData)
	)

	// marshal plan value into data type, appending errors to response
	resp.Diagnostics.Append(req.State.Get(ctx, planData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// extract role name
	roleName = planData.RoleName.ValueString()

	// execute delete call
	{
		osReq := &client.PluginSecurityRoleDeleteRequest{
			Name: roleName,
		}

		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		osResp, err := osReq.Do(ctx, r.client)
		if err != nil {
			if m, ok := err.(client.APIStatusResponse); ok {
				m.AppendDiagnostics(resp.Diagnostics)
			} else {
				resp.Diagnostics.AddError(
					"Error deleting role",
					fmt.Sprintf("Error occurred deleting role %q: %v", roleName, err),
				)
			}
			return
		}

		// attempt to parse response
		sink := client.APIStatusResponse{}
		if err = client.ParseResponse(osResp, &sink, http.StatusOK); err != nil {
			if m, ok := err.(client.APIStatusResponse); ok {
				m.AppendDiagnostics(resp.Diagnostics)
			} else {
				resp.Diagnostics.AddError(
					"Error parsing create role response",
					err.Error(),
				)
			}
			return
		}

		if sink.HasErrors() {
			sink.AppendDiagnostics(resp.Diagnostics)
			return
		}
	}
}

func (r *PluginSecurityRoleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var (
		roleName    string
		osRole      client.PluginSecurityRole
		updateDiags diag.Diagnostics
		ok          bool

		stateData = new(PluginSecurityRoleResourceData)
	)

	// extract role name
	roleName = req.ID

	// query for role from cluster
	// done in sub-context to avoid poisoning ctx var
	{
		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		_, osRoles, err := tryFetchRoles(ctx, r.client, roleName)
		if err != nil {
			if m, ok := err.(client.APIStatusResponse); ok {
				m.AppendDiagnostics(resp.Diagnostics)
			} else {
				resp.Diagnostics.AddError(
					"Error querying for role",
					fmt.Sprintf("Error occurred querying for role %q: %v", roleName, err.Error()),
				)
			}
			return
		}

		// attempt to extract role from response
		if osRole, ok = osRoles[roleName]; !ok {
			resp.Diagnostics.AddError(
				"Role not found",
				fmt.Sprintf("Role %q not found", roleName),
			)
			return
		}
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
