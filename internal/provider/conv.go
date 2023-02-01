package provider

import (
	"github.com/dcarbone/terraform-plugin-framework-utils/v3/conv"
	"github.com/dcarbone/terraform-provider-opensearch/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type (
	attrTypeMap map[string]attr.Type
)

var (
	indexPermissionAttrTypeMap = attrTypeMap{
		resourceAttrIndexPatterns:  types.ListType{ElemType: types.StringType},
		resourceAttrDLS:            types.StringType,
		resourceAttrFLS:            types.StringType,
		resourceAttrMaskedFields:   types.ListType{ElemType: types.StringType},
		resourceAttrAllowedActions: types.ListType{ElemType: types.StringType},
	}

	tenantPermissionAttrTypeMap = attrTypeMap{
		resourceAttrTenantPatterns: types.ListType{ElemType: types.StringType},
		resourceAttrAllowedActions: types.ListType{ElemType: types.StringType},
	}
)

func toNestedObjectList[T any](attrTypes attrTypeMap, in []T, nullOnEmpty bool, fn func(T) (types.Object, diag.Diagnostics)) (types.List, diag.Diagnostics) {
	inLen := len(in)
	objType := types.ObjectType{AttrTypes: attrTypes}

	if nullOnEmpty && inLen == 0 {
		return types.ListNull(objType), nil
	}

	elems := make([]attr.Value, inLen)
	for i, n := range in {
		obj, diags := fn(n)
		if diags.HasError() {
			return types.ListNull(objType), diags
		}
		elems[i] = obj
	}

	return types.ListValue(objType, elems)
}

func indexPermissionToTerraformObject(p client.PluginSecurityRoleIndexPermission) (types.Object, diag.Diagnostics) {
	return types.ObjectValue(
		indexPermissionAttrTypeMap,
		map[string]attr.Value{
			resourceAttrIndexPatterns:  conv.StringsToStringList(p.IndexPatterns, false),
			resourceAttrDLS:            types.StringValue(p.DLS),
			resourceAttrFLS:            types.StringValue(p.FLS),
			resourceAttrMaskedFields:   conv.StringsToStringList(p.MaskedFields, false),
			resourceAttrAllowedActions: conv.StringsToStringList(p.AllowedActions, false),
		},
	)
}

func indexPermissionsToTerraformNestedList(ip []client.PluginSecurityRoleIndexPermission, nullOnEmpty bool) (types.List, diag.Diagnostics) {
	return toNestedObjectList(indexPermissionAttrTypeMap, ip, nullOnEmpty, indexPermissionToTerraformObject)
}

func tenantPermissionToTerraformObject(p client.PluginSecurityRoleTenantPermission) (types.Object, diag.Diagnostics) {
	return types.ObjectValue(
		tenantPermissionAttrTypeMap,
		map[string]attr.Value{
			resourceAttrTenantPatterns: conv.StringsToStringList(p.TenantPatterns, false),
			resourceAttrAllowedActions: conv.StringsToStringList(p.AllowedActions, false),
		},
	)
}

func tenantPermissionsToTerraformNestedList(tp []client.PluginSecurityRoleTenantPermission, nullOnEmpty bool) (types.List, diag.Diagnostics) {
	return toNestedObjectList(tenantPermissionAttrTypeMap, tp, nullOnEmpty, tenantPermissionToTerraformObject)
}

func mapObjectToType[T any](obj types.Object, fn func(map[string]attr.Value) T) T {
	return fn(obj.Attributes())
}

func mapNestedListObjectsToTypes[T any](list types.List, fn func(map[string]attr.Value) T) []T {
	// get all elements in list
	elems := list.Elements()
	elemLen := len(elems)

	out := make([]T, elemLen)
	if elemLen == 0 {
		return out
	}

	for i, e := range elems {
		// this will cause a panic if you didn't do it right.
		attrs := e.(types.Object).Attributes()

		out[i] = fn(attrs)
	}

	return out
}

func mapTerraformIndexPermissionToIndexPermissionType(attrs map[string]attr.Value) client.PluginSecurityRoleIndexPermission {
	// create instance
	out := client.PluginSecurityRoleIndexPermission{}

	// populate
	if v, ok := attrs[resourceAttrIndexPatterns]; ok {
		out.IndexPatterns = conv.StringListToStrings(v)
	}
	if v, ok := attrs[resourceAttrDLS]; ok {
		out.DLS = v.(types.String).ValueString()
	}
	if v, ok := attrs[resourceAttrFLS]; ok {
		out.FLS = v.(types.String).ValueString()
	}
	if v, ok := attrs[resourceAttrMaskedFields]; ok {
		out.MaskedFields = conv.StringListToStrings(v)
	}
	if v, ok := attrs[resourceAttrAllowedActions]; ok {
		out.AllowedActions = conv.StringListToStrings(v)
	}

	// return populated instance
	return out
}

func mapTerraformTenantPermissionsToTenantPermissionsType(attrs map[string]attr.Value) client.PluginSecurityRoleTenantPermission {
	// create instance
	out := client.PluginSecurityRoleTenantPermission{}

	if v, ok := attrs[resourceAttrTenantPatterns]; ok {
		out.TenantPatterns = conv.StringListToStrings(v)
	}
	if v, ok := attrs[resourceAttrAllowedActions]; ok {
		out.AllowedActions = conv.StringListToStrings(v)
	}

	// return populated instance
	return out
}

func terraformSecurityRoleToSecurityRole(d *PluginSecurityRoleResourceData) client.PluginSecurityRole {
	osRole := client.PluginSecurityRole{
		RoleName:    d.RoleName.ValueString(),
		Description: d.Description.ValueString(),

		ClusterPermissions: conv.StringListToStrings(d.ClusterPermissions),
		IndexPermissions:   mapNestedListObjectsToTypes(d.IndexPermissions, mapTerraformIndexPermissionToIndexPermissionType),
		TenantPermissions:  mapNestedListObjectsToTypes(d.TenantPermissions, mapTerraformTenantPermissionsToTenantPermissionsType),
	}

	return osRole
}
