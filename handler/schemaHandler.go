package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/qri-io/jsonpointer"
	"github.com/qri-io/jsonschema"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type SchemaHandler interface {
	LoadSchema(jsonSchemaPath *string) error
	Validate(data *map[string]interface{}) *map[string]string
	AssignControlParameters(c *Control)
	GetValueTypeAtScope(scope *string) string
}

type QriIoSchemaHandler struct {
	Schema *jsonschema.Schema
}

func (h *QriIoSchemaHandler) LoadSchema(jsonSchemaPath *string) error {
	content, err := os.ReadFile(*jsonSchemaPath)
	if err != nil {
		return errors.New(fmt.Sprintf("Error when opening qri-io json schema file: %v ", err.Error()))
	}
	schemaRs := &jsonschema.Schema{}
	if err := json.Unmarshal(content, schemaRs); err != nil {
		return errors.New(fmt.Sprintf("Not valid qri-io json at: %v: %v ", jsonSchemaPath, err.Error()))
	}

	h.Schema = schemaRs

	return nil
}

func (h *QriIoSchemaHandler) Validate(data *map[string]interface{}) *map[string]string {
	ctx := context.Background()
	errorMap := make(map[string]string)

	validationErrors := h.Schema.Validate(ctx, *data)
	for _, vE := range *validationErrors.Errs {
		scopePath := PropertyPathToScopePathPointer(&vE.PropertyPath)
		// Temporary message Length cap
		messageLength := min(len(vE.Message), 40)
		errorMap[strings.Join(scopePath, "/")] = vE.Message[:messageLength]
	}

	return &errorMap

}

// TODO - prefetch type for a #/parameters/field1 like path from a JsonSchema
func (h *QriIoSchemaHandler) GetValueTypeAtScope(scope *string) string {
	return ""
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func (h *QriIoSchemaHandler) AssignControlParameters(c *Control) {

	scopePath, _ := strings.CutPrefix(c.Scope, "#")
	scopePath, _ = strings.CutPrefix(scopePath, "/")

	scopePointerStrings := strings.Split(scopePath, "/")
	pointer := jsonpointer.Pointer{}
	pointer = append(pointer, scopePointerStrings...)

	scopeSchema := h.Schema.Resolve(pointer, scopePath)

	parentSchema := h.Schema.Resolve(pointer[:len(pointer)-2], "") // also remove "properties"- later search parent with "object type"

	requiredControls := []string{}

	if parentSchema != nil {
		reqs := parentSchema.JSONProp("required")

		// fmt.Println("=---", reflect.TypeOf(reqs))
		switch reqs.(type) {
		case *jsonschema.Required:
			reqsList, _ := reqs.(*jsonschema.Required)
			requiredControls = *reqsList
		}
	}

	if scopeSchema != nil {
		c.InputType = fmt.Sprint(scopeSchema.JSONProp("type"))

		if c.InputType == "integer" {
			c.InputType = "number"
		} else if c.InputType == "string" {
			c.InputType = "text"

			rawFormat := scopeSchema.JSONProp("format")
			// fmt.Println("=---", rawFormat, reflect.TypeOf(rawFormat))
			switch rawFormat.(type) {
			case *jsonschema.Format:
				val := rawFormat.(*jsonschema.Format)
				if *val == "date" {
					c.InputType = "date"
				}
				break
			}

			enums := scopeSchema.JSONProp("enum")

			switch enums.(type) {
			case *jsonschema.Enum:
				enumOptions := enums.(*jsonschema.Enum)
				optionsAsStrings := []string{}
				for _, o := range *enumOptions {
					optionsAsStrings = append(optionsAsStrings, o.String())
				}
				c.EnumOptions = optionsAsStrings
				break
			}
		}

		description := scopeSchema.JSONProp("description")
		switch description.(type) {
		case *jsonschema.Description:
			c.Description = string(*description.(*jsonschema.Description))
			break

		}
	}

	labelLowerCase := lastString(strings.Split(scopePath, "/"))
	c.Label = cases.Title(language.Und, cases.Compact).String(labelLowerCase)

	c.Required = Contains(len(requiredControls), func(i int) bool {
		return requiredControls[i] == labelLowerCase
	})
}

func lastString(ss []string) string {
	return ss[len(ss)-1]
}
