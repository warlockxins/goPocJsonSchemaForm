package main

import (
	"got/handler"
	"testing"
)

func TestHelloEmpty(t *testing.T) {

	layoutHandler := handler.LayoutHandler{}
	res, err := layoutHandler.GetSchemaAndBindLayout("./screens/contact/uiSchema.json", "./screens/contact/jsonSchema.json")

	if res == nil || err != nil {
		t.Fatalf(`not found or unparsable %v`, err)
	}

	if res.Type != "VerticalLayout" {
		t.Fatalf("root not VerticalLayout")
	}

	control := res.VerticalLayout.Elements[0].HorizontalLayout.Elements[0].Control
	if control.Label != "Name" {
		t.Fatalf("Control label should be 'Name'")
	}

	if control.Description != "Please enter your name" {
		t.Fatalf("Control label should be 'Name'")
	}

	nationality := res.VerticalLayout.Elements[2].HorizontalLayout.Elements[1]

	if len(nationality.Control.EnumOptions) != 6 {
		t.Fatalf("nationality should have 6 enum option entries")
	}
}
