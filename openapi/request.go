package openapi

// Param ParameterObject
type Param struct {
	root *OpenAPI
	// Fixed fields
	Name            string    `json:"name" validate:"required"`
	In              ParamType `json:"in" validate:"required,oneof=query header path cookie"`
	Description     string    `json:"description,omitempty"`
	Required        bool      `json:"required"`
	Deprecated      bool      `json:"deprecated,omitempty"`
	AllowEmptyValue bool      `json:"allowEmptyValue,omitempty"`
	// Below are optional fields
	Schema   *Schema             `json:"schema,omitempty"`
	Example  interface{}         `json:"example,omitempty"`
	Examples map[string]*Example `json:"examples,omitempty"`
}

// RequestBody request body object
type RequestBody struct {
	Description string `json:"description,omitempty"`
	// MIME-Type -> MediaTypeObject
	Content  mediaTypeMap `json:"content" validate:"required"`
	Required bool         `json:"required,omitempty"`
}
