package openapi

import (
	"reflect"
	"strconv"
)

// Operation OperationObject
type Operation struct {
	method       string
	path         *Path
	MTags        []string     `json:"tags,omitempty"`
	MSummary     string       `json:"summary,omitempty"`
	MDescription string       `json:"description,omitempty"`
	MOperationID string       `json:"operationId,omitempty"`
	MParameters  []*Param     `json:"parameters,omitempty"`
	MRequestBody *RequestBody `json:"requestBody,omitempty"`
	// HTTP Status Code => Response
	MResponses  Responses `json:"responses" validate:"required"`
	MDeprecated bool      `json:"deprecated,omitempty"`
}

// Summary setter
func (o *Operation) Summary(s string) *Operation {
	o.MSummary = s
	return o
}

// Description setter
func (o *Operation) Description(s string) *Operation {
	o.MDescription = s
	return o
}

// Tags setter
func (o *Operation) Tags(tags ...string) *Operation {
	tagMap := make(map[string]struct{})
	for _, tag := range o.MTags {
		tagMap[tag] = struct{}{}
	}
	for _, tag := range tags {
		if _, ok := tagMap[tag]; !ok {
			tagMap[tag] = struct{}{}
			o.MTags = append(o.MTags, tag)
		}
	}
	return o
}

// Param add a parameter
func (o *Operation) Param(p *Param) *Operation {
	for _, param := range o.MParameters {
		if param.Name == p.Name {
			panic("duplicate parameter of name:" + p.Name)
		}
	}
	paramID := o.MOperationID + "_" + p.Name
	o.Root().Components.Parameters[paramID] = p
	return o
}

// Metadata add metadata to operation.
// OperationID is required in OpenAPI. If it's empty, an operation id is generated automatically
func (o *Operation) Metadata(summary, description string) *Operation {
	o.MSummary = summary
	o.MDescription = description
	return o
}

// RequestBody setup request body
func (o *Operation) RequestBody(fn func(*RequestBody)) *Operation {
	if o.MRequestBody == nil {
		o.MRequestBody = new(RequestBody)
	}
	fn(o.MRequestBody)
	o.Root().Components.RequestBodies[o.MOperationID] = o.MRequestBody
	return o
}

// InJSON request body in json format
func (o *Operation) InJSON(v interface{}, bodyFns ...func(*RequestBody)) *Operation {
	rb := &RequestBody{}
	o.Root().Components.RequestBodies[o.MOperationID] = rb
	for _, fn := range bodyFns {
		fn(rb)
	}
	return o
}

// OutJSON response in json format
func (o *Operation) OutJSON(code int, v interface{}, bodyFn ...func(*Response)) *Operation {
	strCode := strconv.Itoa(code)
	if _, exists := o.MResponses[strCode]; exists {
		panic("operation " + o.MOperationID + " already returns code " + strCode)
	}
	schema := o.Root().MustGetSchema("", v)
	resp := &Response{
		Headers: make(paramMap),
		Content: mediaTypeMap{
			MimeJSON: &MediaType{
				Schema:  schema,
				Example: v,
			},
		},
	}
	o.MResponses[strCode] = resp
	for _, fn := range bodyFn {
		fn(resp)
	}
	return o
}

// InForm add a param in form
func (o *Operation) InForm(v interface{}, bodyFn ...func(*Response)) *Operation {
	// TODO: parse interface with "form" tag
	schema := o.Root().MustGetSchema("", v)
	resp := &Response{
		Headers: make(paramMap),
		Content: mediaTypeMap{
			MimeJSON: &MediaType{
				Schema:  schema,
				Example: v,
			},
		},
	}
	for _, fn := range bodyFn {
		fn(resp)
	}
	return o
}

// Root returns document root for operation
func (o *Operation) Root() *OpenAPI {
	if o.path == nil {
		panic("no path is set for operation")
	}
	return o.path.Root()
}

// Returns with code
func (o *Operation) Returns(code int, description string, key string, v interface{}) *Operation {
	strCode := strconv.Itoa(code)
	if _, exists := o.MResponses[strCode]; exists {
		panic("operation " + o.MOperationID + " already returns code " + strCode)
	}
	r := o.newResponse(description, key, v)
	o.MResponses[strCode] = r
	return o
}

// ReturnsNonJSON return something not json
func (o *Operation) ReturnsNonJSON(code int, description string,
	mimeType string, headers map[string]*Param, schema *Schema, example interface{}) *Operation {
	strCode := strconv.Itoa(code)
	if _, exists := o.MResponses[strCode]; exists {
		panic("operation " + o.MOperationID + " already returns code " + strCode)
	}
	o.MResponses[strCode] = &Response{
		Description: description,
		Headers:     headers,
		Content: mediaTypeMap{
			mimeType: &MediaType{
				Schema:  schema,
				Example: example,
			},
		},
	}
	return o
}

// ReturnDefault add default response.
// A default response is the response to be used when none of defined codes match the situation.
func (o *Operation) ReturnDefault(description string, key string, v interface{}) *Operation {
	o.MResponses["default"] = o.newResponse(description, key, v)
	return o
}

func (o *Operation) newResponse(description string, key string, v interface{}) *Response {
	schema := o.Root().MustGetSchema(key, v)
	return &Response{
		Description: description,
		Headers:     make(paramMap),
		Content: mediaTypeMap{
			MimeJSON: &MediaType{
				Schema:  schema,
				Example: v,
			},
		},
	}
}

// ReadJSON read object json from request body
func (o *Operation) ReadJSON(description string, required bool, key string, v interface{}) *Operation {
	schema := o.Root().MustGetSchema(key, v)
	o.MRequestBody = &RequestBody{
		Description: description,
		Required:    required,
		Content: mediaTypeMap{
			MimeJSON: &MediaType{
				Schema:  schema,
				Example: v,
			},
		},
	}
	return o
}

// Read read raw body of any kind
func (o *Operation) Read(description string, required bool, mimeType string, example interface{}) *Operation {
	o.MRequestBody = &RequestBody{
		Description: description,
		Required:    required,
		Content: mediaTypeMap{
			mimeType: &MediaType{
				Example: example,
			},
		},
	}
	return o
}

// AddParam with param in operation
func (o *Operation) AddParam(in ParamType, name, description string) *Param {
	if !in.IsValid() {
		panic("invalid param in " + in)
	}
	param := &Param{
		In:          in,
		Name:        name,
		Description: description,
	}
	// A path param is always required
	if in == PathParam {
		param.Required = true
		param.Schema = &Schema{
			Type: "string",
		}
	}
	o.MParameters = append(o.MParameters, param)
	return param
}

// WithParam add param to operation
func (o *Operation) WithParam(param *Param) *Operation {
	if !param.In.IsValid() {
		panic("invalid param in " + param.In)
	}
	o.MParameters = append(o.MParameters, param)
	return o
}

// WithPathParam add path param
func (o *Operation) WithPathParam(name, description string) *Operation {
	return o.WithParam(&Param{
		In:          PathParam,
		Name:        name,
		Description: description,
		Required:    true,
		Schema: &Schema{
			Type: "string",
		},
	})
}

// WithQueryParam add query param. Complex types of query param is not supported here(e.g., a struct or slice)
func (o *Operation) WithQueryParam(name, description string, example interface{}) *Operation {
	tv := reflect.TypeOf(example)
	typ, _ := kindToType(tv.Kind())
	return o.WithParam(&Param{
		In:          QueryParam,
		Name:        name,
		Description: description,
		Example:     example,
		Schema: &Schema{
			Type: typ,
		},
	})
}
