package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

type LayoutHandler struct{}

type Screen struct {
	Name    string
	Current bool
}

type LayoutPage struct {
	Screens           []Screen
	Layout            *Layout
	CurrentScreenPost string
}

func (h *LayoutHandler) HandleLayoutShow(c echo.Context, layoutPage *LayoutPage) error {
	return c.Render(http.StatusOK, "formIndex", layoutPage)
}

func (h *LayoutHandler) HandleFormShow(c echo.Context, layout *Layout) error {
	return c.Render(http.StatusOK, "formTypeSwitch", layout)
}

func (h *LayoutHandler) GetUiLayout(uiSchemaPath string) (*Layout, error) {
	content, err := os.ReadFile(uiSchemaPath)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error when opening file: %v ", err.Error()))
	}
	res := Layout{}
	err = json.Unmarshal([]byte(content), &res)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Not valid json at: %v: %v ", uiSchemaPath, err.Error()))
	}

	return &res, nil
}
