package handler

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type FormData struct {
	Values url.Values
	Errors map[string]string
}

func NewFormData() *FormData {
	return &FormData{
		Values: make(url.Values),
		Errors: make(map[string]string),
	}
}

type TypeSwitch struct {
	Type string `json:"type"`
}

type Layout struct {
	TypeSwitch
	*VerticalLayout
	*HorizontalLayout
	*Control
	*Label
}

type Label struct {
	Text string `json:"text"`
}

type VerticalLayout struct {
	Elements []Layout `json:"elements"`
}

type HorizontalLayout struct {
	Elements []Layout `json:"elements"`
}

type ControlBase struct {
	Suggestion  []string `json:"suggestion"`
	InputType   string
	Label       string
	EnumOptions []string
	Required    bool
	Description string `json:"description"`
}

type Control struct {
	ControlBase
	Scope string      `json:"scope"`
	Value interface{} // can be int, string, null?
	Error interface{}
}

func (c *Control) ExtractTypeAndNameValueError(sH *SchemaHandler, existingValues *FormData) {
	c.Value = existingValues.Values.Get(c.Scope)
	c.Error = existingValues.Errors[c.Scope]

	(*sH).AssignControlParameters(c)
}

func Contains(n int, match func(i int) bool) bool {
	for i := 0; i < n; i++ {
		if match(i) {
			return true
		}
	}
	return false
}

func (t *Layout) BindToSchema(sH *SchemaHandler, existingValues *FormData) {
	switch t.Type {
	case "VerticalLayout":
		for _, e := range t.VerticalLayout.Elements {
			e.BindToSchema(sH, existingValues)
		}
		return
	case "HorizontalLayout":
		for _, e := range t.HorizontalLayout.Elements {
			e.BindToSchema(sH, existingValues)
		}
		return
	case "Control":
		t.Control.ExtractTypeAndNameValueError(sH, existingValues)
		return
	}

}

func (t *Layout) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &t.TypeSwitch); err != nil {
		return err
	}

	switch t.Type {
	case "VerticalLayout":
		t.VerticalLayout = &VerticalLayout{}
		return json.Unmarshal(data, t.VerticalLayout)
	case "HorizontalLayout":
		t.HorizontalLayout = &HorizontalLayout{}
		return json.Unmarshal(data, t.HorizontalLayout)

	case "Control":
		t.Control = &Control{}
		return json.Unmarshal(data, t.Control)

	case "Label":
		t.Label = &Label{}
		return json.Unmarshal(data, t.Label)
	default:
		return fmt.Errorf("unrecognized type value %q", t.Type)
	}

}
