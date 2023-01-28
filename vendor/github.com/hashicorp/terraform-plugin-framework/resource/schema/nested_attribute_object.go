package schema

import (
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema/fwxschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// Ensure the implementation satisifies the desired interfaces.
var (
	_ fwxschema.NestedAttributeObjectWithPlanModifiers = NestedAttributeObject{}
	_ fwxschema.NestedAttributeObjectWithValidators    = NestedAttributeObject{}
)

// NestedAttributeObject is the object containing the underlying attributes
// for a ListNestedAttribute, MapNestedAttribute, SetNestedAttribute, or
// SingleNestedAttribute (automatically generated). When retrieving the value
// for this attribute, use types.Object as the value type unless the CustomType
// field is set. The Attributes field must be set. Nested attributes are only
// compatible with protocol version 6.
//
// This object enables customizing and simplifying details within its parent
// NestedAttribute, therefore it cannot have Terraform schema fields such as
// Required, Description, etc.
type NestedAttributeObject struct {
	// Attributes is the mapping of underlying attribute names to attribute
	// definitions. This field must be set.
	Attributes map[string]Attribute

	// CustomType enables the use of a custom attribute type in place of the
	// default basetypes.ObjectType. When retrieving data, the basetypes.ObjectValuable
	// associated with this custom type must be used in place of types.Object.
	CustomType basetypes.ObjectTypable

	// Validators define value validation functionality for the attribute. All
	// elements of the slice of AttributeValidator are run, regardless of any
	// previous error diagnostics.
	//
	// Many common use case validators can be found in the
	// github.com/hashicorp/terraform-plugin-framework-validators Go module.
	//
	// If the Type field points to a custom type that implements the
	// xattr.TypeWithValidate interface, the validators defined in this field
	// are run in addition to the validation defined by the type.
	Validators []validator.Object

	// PlanModifiers defines a sequence of modifiers for this attribute at
	// plan time. Schema-based plan modifications occur before any
	// resource-level plan modifications.
	//
	// Schema-based plan modifications can adjust Terraform's plan by:
	//
	//  - Requiring resource recreation. Typically used for configuration
	//    updates which cannot be done in-place.
	//  - Setting the planned value. Typically used for enhancing the plan
	//    to replace unknown values. Computed must be true or Terraform will
	//    return an error. If the plan value is known due to a known
	//    configuration value, the plan value cannot be changed or Terraform
	//    will return an error.
	//
	// Any errors will prevent further execution of this sequence or modifiers.
	PlanModifiers []planmodifier.Object
}

// ApplyTerraform5AttributePathStep performs an AttributeName step on the
// underlying attributes or returns an error.
func (o NestedAttributeObject) ApplyTerraform5AttributePathStep(step tftypes.AttributePathStep) (any, error) {
	return fwschema.NestedAttributeObjectApplyTerraform5AttributePathStep(o, step)
}

// Equal returns true if the given NestedAttributeObject is equivalent.
func (o NestedAttributeObject) Equal(other fwschema.NestedAttributeObject) bool {
	if _, ok := other.(NestedAttributeObject); !ok {
		return false
	}

	return fwschema.NestedAttributeObjectEqual(o, other)
}

// GetAttributes returns the Attributes field value.
func (o NestedAttributeObject) GetAttributes() fwschema.UnderlyingAttributes {
	return schemaAttributes(o.Attributes)
}

// ObjectPlanModifiers returns the PlanModifiers field value.
func (o NestedAttributeObject) ObjectPlanModifiers() []planmodifier.Object {
	return o.PlanModifiers
}

// ObjectValidators returns the Validators field value.
func (o NestedAttributeObject) ObjectValidators() []validator.Object {
	return o.Validators
}

// Type returns the framework type of the NestedAttributeObject.
func (o NestedAttributeObject) Type() basetypes.ObjectTypable {
	if o.CustomType != nil {
		return o.CustomType
	}

	return fwschema.NestedAttributeObjectType(o)
}
