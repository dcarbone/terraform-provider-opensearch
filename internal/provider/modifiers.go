package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type pluginSecurityRoleIndexSettingsDefaultValue struct {
}

func (m *pluginSecurityRoleIndexSettingsDefaultValue) Description(context.Context) string {
	return "Ensures state consistency when no value is provided in state"
}

func (m *pluginSecurityRoleIndexSettingsDefaultValue) MarkdownDescription(context.Context) string {
	return "Ensures state consistency when no value is provided in state"
}

func (m *pluginSecurityRoleIndexSettingsDefaultValue) PlanModifyList(_ context.Context, req planmodifier.ListRequest, resp *planmodifier.ListResponse) {
	if req.PlanValue.IsNull() {
		resp.PlanValue = types.ListValueMust(types.ObjectType{AttrTypes: indexPermissionAttrTypeMap}, []attr.Value{})
	}
}

type pluginSecurityRoleTenantPermissionsDefaultValue struct {
}

func (m *pluginSecurityRoleTenantPermissionsDefaultValue) Description(context.Context) string {
	return "Ensures state consistency when no value is provided in state"
}

func (m *pluginSecurityRoleTenantPermissionsDefaultValue) MarkdownDescription(context.Context) string {
	return "Ensures state consistency when no value is provided in state"
}

func (m *pluginSecurityRoleTenantPermissionsDefaultValue) PlanModifyList(_ context.Context, req planmodifier.ListRequest, resp *planmodifier.ListResponse) {
	if req.PlanValue.IsNull() {
		resp.PlanValue = types.ListValueMust(types.ObjectType{AttrTypes: tenantPermissionAttrTypeMap}, []attr.Value{})
	}
}

type defaultValuedStringPlanModifier string

func (defaultValuedStringPlanModifier) Description(context.Context) string {
	const d = "Sets a default value when state is null"
	return d
}

func (defaultValuedStringPlanModifier) MarkdownDescription(context.Context) string {
	const d = "Sets a default value when state is null"
	return d
}

func (m defaultValuedStringPlanModifier) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.PlanValue.IsNull() {
		resp.PlanValue = types.StringValue(string(m))
	}
}
