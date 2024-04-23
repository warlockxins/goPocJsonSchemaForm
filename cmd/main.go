package main

import (
	"fmt"
	"got/handler"
	"html"
	"io"
	"path"
	"text/template"

	"github.com/labstack/echo/v4"
)

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func newTeplate() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
}

func newQriIoJsonSchemaHandler() handler.SchemaHandler {
	return &handler.QriIoSchemaHandler{
		Schema: nil,
	}
}

func main() {

	layoutHandler := handler.LayoutHandler{}

	app := echo.New()

	app.Renderer = newTeplate()

	app.Static("/images", "images")
	app.Static("/css", "css")

	app.GET("/", func(c echo.Context) error {
		directories, err := handler.GetScreens()
		if err != nil {
			return c.String(500, "cannot get screen list")
		}

		return layoutHandler.HandleLayoutShow(c, &handler.LayoutPage{
			Screens: handler.ScreensToHandlerScreeenList(&directories, nil),
			Layout:  nil,
		})

	})

	app.GET("/layout/:id", func(c echo.Context) error {

		currentPage := c.Param("id")

		layoutPath := fmt.Sprintf("./screens/%s/uiSchema.json", currentPage)
		schemaPath := fmt.Sprintf("./screens/%s/jsonSchema.json", currentPage)

		jsonSchemaHandler := newQriIoJsonSchemaHandler()
		err := jsonSchemaHandler.LoadSchema(&schemaPath)

		if err != nil {
			fmt.Println("Error getting data schema:", err)
			return c.String(500, "Incorrect data schema")
		}

		res, err := layoutHandler.GetUiLayout(layoutPath)

		if err != nil {
			fmt.Println("Error getting layout schema:", err)
			return c.String(500, "Incorrect Ui schema")
		}

		directories, err := handler.GetScreens()

		if err != nil {
			return c.String(500, "cannot get screen list")
		}

		res.BindToSchema(&jsonSchemaHandler, handler.NewFormData())

		return layoutHandler.HandleLayoutShow(c, &handler.LayoutPage{
			Screens:           handler.ScreensToHandlerScreeenList(&directories, &currentPage),
			Layout:            res,
			CurrentScreenPost: html.EscapeString(path.Join("/layout", currentPage)),
		})
	})

	app.POST("/layout/:id", func(c echo.Context) error {

		formParams, formParamsErrors := c.FormParams()
		if formParamsErrors != nil {
			fmt.Println("form data cannot be read", formParamsErrors)
			return c.String(500, "Problem getting form data")
		}

		currentPage := c.Param("id")

		layoutPath := fmt.Sprintf("./screens/%s/uiSchema.json", currentPage)
		schemaPath := fmt.Sprintf("./screens/%s/jsonSchema.json", currentPage)

		jsonSchemaHandler := newQriIoJsonSchemaHandler()
		err := jsonSchemaHandler.LoadSchema(&schemaPath)

		if err != nil {
			fmt.Println("Error getting data schema:", err)
			return c.String(500, "Incorrect data schema")
		}

		res, err := layoutHandler.GetUiLayout(layoutPath)

		if err != nil {
			fmt.Println("Error getting layout schema:", err)
			return c.String(500, "Incorrect Ui schema")
		}

		receivedValueMapForValidation := handler.FormValuesToGroupedMap(&formParams, &jsonSchemaHandler)
		fmt.Println("here is deep obj", receivedValueMapForValidation)

		formData := handler.FormData{
			Values: formParams,
			Errors: *jsonSchemaHandler.Validate(&receivedValueMapForValidation),
		}

		res.BindToSchema(&jsonSchemaHandler, &formData)

		return layoutHandler.HandleFormShow(c, res)
	})

	app.Start(":3000")

}
