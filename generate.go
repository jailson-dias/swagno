package swagno

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
)

var ErrIgnoreJSONField = errors.New("ignore json field")

// GenerateDocs Generate swagger v2 documentation as json string
func (s Swagger) GenerateDocs() (jsonDocs []byte) {
	if len(endpoints) == 0 {
		log.Println("No endpoints found")
		return
	}

	// generate definition object of s json: https://swagger.io/specification/v2/#definitions-object
	s.generateSwaggerDefinition(endpoints)

	// convert all user EndPoint models to 'path' fields of s json
	// https://swagger.io/specification/v2/#paths-object
	for _, endpoint := range endpoints {
		path := endpoint.Path

		if s.Paths[path] == nil {
			s.Paths[path] = make(map[string]swaggerEndpoint)
		}

		method := strings.ToLower(endpoint.Method)

		consumes := []string{ContentTypeApplicationJSON}
		produces := []string{ContentTypeApplicationJSON, ContentTypeApplicationXML}
		for _, param := range endpoint.Params {
			if param.In == "formData" {
				consumes = append([]string{"multipart/form-data"}, consumes...)
				break
			}
		}
		if len(endpoint.Consume) == 0 {
			consumes = append(endpoint.Consume, consumes...)
		}
		if len(endpoint.Produce) == 0 {
			produces = append(endpoint.Produce, produces...)
		}

		parameters := make([]swaggerParameter, 0)
		for _, param := range endpoint.Params {
			parameters = append(parameters, swaggerParameter{
				Name:              param.Name,
				In:                param.In,
				Description:       param.Description,
				Required:          param.Required,
				Type:              param.Type,
				Format:            param.Format,
				Items:             param.Items,
				Enum:              param.Enum,
				Default:           param.Default,
				Min:               param.Min,
				Max:               param.Max,
				MinLen:            param.MinLen,
				MaxLen:            param.MaxLen,
				Pattern:           param.Pattern,
				MaxItems:          param.MaxItems,
				MinItems:          param.MinItems,
				UniqueItems:       param.UniqueItems,
				MultipleOf:        param.MultipleOf,
				CollenctionFormat: param.CollectionFormat,
			})
		}

		if endpoint.Body != nil {
			bodySchema := swaggerResponseScheme{
				Ref: fmt.Sprintf("#/definitions/%T", endpoint.Body),
			}
			if reflect.TypeOf(endpoint.Body).Kind() == reflect.Slice {
				bodySchema = swaggerResponseScheme{
					Type: "array",
					Items: &swaggerResponseSchemeItems{
						Ref: fmt.Sprintf("#/definitions/%T", endpoint.Body),
					},
				}
			}
			parameters = append(parameters, swaggerParameter{
				Name:        "body",
				In:          "body",
				Description: "body",
				Required:    true,
				Schema:      &bodySchema,
			})
		}

		// add each endpoint to paths field of s
		s.Paths[path][method] = swaggerEndpoint{
			Description: endpoint.Description,
			Summary:     endpoint.Description,
			OperationId: method + "-" + path,
			Consumes:    consumes,
			Produces:    produces,
			Tags:        endpoint.Tags,
			Parameters:  parameters,
			Responses:   buildSwaggerResponses(endpoint.Responses),
			Security:    endpoint.Security,
		}
	}

	// convert Swagger instance to json string and return it
	json, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		log.Println("Error while generating s json")
	}
	return json
}

func buildSwaggerResponses(list []Response) map[string]swaggerResponse {
	responses := make(map[string]swaggerResponse)
	for _, response := range list {
		responseSchema := &swaggerResponseScheme{
			Ref: fmt.Sprintf("#/definitions/%T", response.Body),
		}
		if reflect.TypeOf(response.Body).Kind() == reflect.Slice {
			responseSchema = &swaggerResponseScheme{
				Type: "array",
				Items: &swaggerResponseSchemeItems{
					Ref: fmt.Sprintf("#/definitions/%T", response.Body),
				},
			}
		}
		responses[response.Code] = swaggerResponse{
			Description: response.Description,
			Schema:      *responseSchema,
		}
	}
	return responses
}

// generate "definitions" keys from endpoints: https://swagger.io/specification/v2/#definitions-object
func (s *Swagger) generateSwaggerDefinition(endpoints []Endpoint) {
	// create all definations for each model used in endpoint
	s.Definitions = make(map[string]swaggerDefinition)
	for _, endpoint := range endpoints {
		if endpoint.Body != nil {
			s.createDefinition(fmt.Sprintf("%T", endpoint.Body), endpoint.Body)
		}
		for _, response := range endpoint.Responses {
			s.createDefinition(fmt.Sprintf("%T", response.Body), response.Body)
		}
	}
}

// generate "definitions" attribute for swagger json
func (s *Swagger) createDefinition(definition string, obj interface{}) {
	reflectReturn := reflect.TypeOf(obj)
	if reflectReturn.Kind() == reflect.Slice {
		reflectReturn = reflectReturn.Elem()
	}

	if _, ok := s.Definitions[definition]; !ok {
		s.Definitions[definition] = swaggerDefinition{
			Type:       "object",
			Properties: make(map[string]swaggerDefinitionProperties),
		}
	}
	for i := 0; i < reflectReturn.NumField(); i++ {
		field := reflectReturn.Field(i)
		fieldType := getType(field.Type.Kind().String())

		// skip for function and channel types
		if fieldType == "func" || fieldType == "chan" {
			continue
		}

		fieldName, err := getJsonTag(field)
		if err != nil {
			continue
		}

		// if item type is array, create defination for array element type
		if fieldType == "array" {
			if field.Type.Elem().Kind() == reflect.Struct {
				definitionName := field.Type.Elem().String()
				s.Definitions[definition].Properties[fieldName] = swaggerDefinitionProperties{
					Type: fieldType,
					Items: &swaggerDefinitionPropertiesItems{
						Ref: fmt.Sprintf("#/definitions/%s", definitionName),
					},
				}

				if _, ok := s.Definitions[definitionName]; !ok {
					definitionObj := reflect.New(field.Type.Elem()).Elem().Interface()
					s.createDefinition(fmt.Sprintf("%T", definitionObj), definitionObj)
				}
			} else {
				s.Definitions[definition].Properties[fieldName] = swaggerDefinitionProperties{
					Type: fieldType,
					Items: &swaggerDefinitionPropertiesItems{
						Type:    getType(field.Type.Elem().Kind().String()),
						Example: getExampleValue(field),
					},
				}
			}
		} else {
			if field.Type.Kind() == reflect.Struct {
				if field.Type.String() == "time.Time" {
					s.Definitions[definition].Properties[fieldName] = swaggerDefinitionProperties{
						Type:    "string",
						Format:  "date-time",
						Example: getExampleValue(field),
					}
				} else if field.Type.String() == "time.Duration" {
					s.Definitions[definition].Properties[fieldName] = swaggerDefinitionProperties{
						Type:    "integer",
						Example: getExampleValue(field),
					}
				} else {
					jsonTag := field.Tag.Get("json")
					if jsonTag == "" {
						s.createDefinition(fmt.Sprintf("%T", obj), reflect.New(field.Type).Elem().Interface())
					} else {
						definitionName := field.Type.String()
						s.Definitions[definition].Properties[fieldName] = swaggerDefinitionProperties{
							Ref: fmt.Sprintf("#/definitions/%s", definitionName),
						}

						if _, ok := s.Definitions[definitionName]; !ok {
							definitionObj := reflect.New(field.Type).Elem().Interface()
							s.createDefinition(fmt.Sprintf("%T", definitionObj), definitionObj)
						}
					}

				}
			} else if field.Type.Kind() == reflect.Pointer {
				if field.Type.Elem().Kind() == reflect.Struct {
					if field.Type.Elem().String() == "time.Time" {
						s.Definitions[definition].Properties[fieldName] = swaggerDefinitionProperties{
							Type:    "string",
							Format:  "date-time",
							Example: getExampleValue(field),
						}
					} else if field.Type.String() == "time.Duration" {
						s.Definitions[definition].Properties[fieldName] = swaggerDefinitionProperties{
							Type:    "integer",
							Example: getExampleValue(field),
						}
					} else {
						definitionName := field.Type.Elem().String()
						s.Definitions[definition].Properties[fieldName] = swaggerDefinitionProperties{
							Ref: fmt.Sprintf("#/definitions/%s", definitionName),
						}

						if _, ok := s.Definitions[definitionName]; !ok {
							definitionObj := reflect.New(field.Type.Elem()).Elem().Interface()
							s.createDefinition(fmt.Sprintf("%T", definitionObj), definitionObj)
						}
					}
				} else {
					s.Definitions[definition].Properties[fieldName] = swaggerDefinitionProperties{
						Type:    getType(field.Type.Elem().Kind().String()),
						Example: getExampleValue(field),
					}
				}
			} else {
				s.Definitions[definition].Properties[fieldName] = swaggerDefinitionProperties{
					Type:    fieldType,
					Example: getExampleValue(field),
				}
			}
		}
	}
}

// get struct json tag as string of a struct field
func getJsonTag(field reflect.StructField) (string, error) {
	jsonTag := field.Tag.Get("json")
	jsonOptions := strings.Split(jsonTag, ",")

	if len(jsonOptions) > 0 {
		if jsonOptions[0] == "-" {
			return "", ErrIgnoreJSONField
		} else if jsonOptions[0] == "" {
			return field.Name, nil
		}

		return jsonOptions[0], nil
	}

	return field.Name, nil
}

func getEmptyExampleValue(fieldType string) interface{} {
	switch fieldType {
	case "integer":
		return 0
	case "boolean":
		return false
	case "number":
		return 0.0
	case "struct":
		return nil
	}

	return fieldType
}

// get example tag as string example value
func getExampleValue(field reflect.StructField) interface{} {
	example := field.Tag.Get("example")
	fieldType := getType(field.Type.Kind().String())
	if example == "" {
		switch fieldType {
		case "ptr":
			return getEmptyExampleValue(getType(field.Type.Elem().Kind().String()))
		default:
			return getEmptyExampleValue(fieldType)
		}
	}

	switch fieldType {
	case "integer":
		number, err := strconv.Atoi(example)
		if err != nil {
			return 0
		}
		return number
	case "boolean":
		if example == "true" {
			return true
		}

		return false
	case "number":
		number, err := strconv.ParseFloat(example, 64)
		if err != nil {
			return 0.0
		}

		return number
	}

	return example
}

// get swagger type from reflection type
// https://swagger.io/specification/v2/#data-types
func getType(t string) string {
	if strings.Contains(strings.ToLower(t), "int") {
		return "integer"
	} else if t == "array" || t == "slice" {
		return "array"
	} else if t == "bool" {
		return "boolean"
	} else if t == "float64" || t == "float32" {
		return "number"
	}
	return t
}
