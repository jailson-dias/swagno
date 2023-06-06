package swagno

// https://swagger.io/specification/v2/#definitionsObject
type swaggerDefinition struct {
	Type       string                                 `json:"type"`
	Properties map[string]swaggerDefinitionProperties `json:"properties"`
	Required   []string                               `json:"required,omitempty"`
}

// https://swagger.io/specification/v2/#schemaObject
type swaggerDefinitionProperties struct {
	Type    string                            `json:"type,omitempty"`
	Format  string                            `json:"format,omitempty"`
	Ref     string                            `json:"$ref,omitempty"`
	Items   *swaggerDefinitionPropertiesItems `json:"items,omitempty"`
	Example interface{}                       `json:"example,omitempty"`
	Enum    interface{}                       `json:"enum,omitempty"`
}

type swaggerDefinitionPropertiesItems struct {
	Type    string      `json:"type,omitempty"`
	Ref     string      `json:"$ref,omitempty"`
	Example interface{} `json:"example,omitempty"`
	Enum    interface{} `json:"enum,omitempty"`
}
