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

	nestedObjectFunc func(any) (types.Object, diag.Diagnostics)
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

func indexPermissionToObject(p client.PluginSecurityRoleIndexPermission) (types.Object, diag.Diagnostics) {
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

func indexPermissionsToNestedList(ip []client.PluginSecurityRoleIndexPermission, nullOnEmpty bool) (types.List, diag.Diagnostics) {
	return toNestedObjectList(indexPermissionAttrTypeMap, ip, nullOnEmpty, indexPermissionToObject)
}

func tenantPermissionToObject(p client.PluginSecurityRoleTenantPermission) (types.Object, diag.Diagnostics) {
	return types.ObjectValue(
		tenantPermissionAttrTypeMap,
		map[string]attr.Value{
			resourceAttrTenantPatterns: conv.StringsToStringList(p.TenantPatterns, false),
			resourceAttrAllowedActions: conv.StringsToStringList(p.AllowedActions, false),
		},
	)
}

func tenantPermissionsToNestedList(tp []client.PluginSecurityRoleTenantPermission, nullOnEmpty bool) (types.List, diag.Diagnostics) {
	return toNestedObjectList(tenantPermissionAttrTypeMap, tp, nullOnEmpty, tenantPermissionToObject)
}
