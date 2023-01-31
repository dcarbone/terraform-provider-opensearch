package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewSecurityPluginRoleResource() resource.Resource {
	r := new(SecurityPluginRoleResource)
	return r
}

type SecurityPluginRoleResource struct {
	ResourceShared
}

type SecurityPluginRoleResourceData struct {
	RoleName           types.String `tfsdk:"role_name"`
	Description        types.String `tfsdk:"description"`
	ClusterPermissions types.List   `tfsdk:"cluster_permissions"`
	IndexPermissions   types.List   `tfsdk:"index_permissions"`
	TenantPermissions  types.List   `tfsdk:"tenant_permissions"`

	Reserved types.Bool `tfsdk:"reserved"`
	Hidden   types.Bool `tfsdk:"hidden"`
	Static   types.Bool `tfsdk:"static"`
}

func (r *SecurityPluginRoleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = makeResourceName(req.ProviderTypeName, resourceSuffixSecurityPluginRole)
}

func (r *SecurityPluginRoleResource) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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

func (r *SecurityPluginRoleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

}

func (r *SecurityPluginRoleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

}

func (r *SecurityPluginRoleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

}

func (r *SecurityPluginRoleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

}

func (r *SecurityPluginRoleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

}
