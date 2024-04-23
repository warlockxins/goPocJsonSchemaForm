package handler

import (
	"net/url"
	"strings"
)

func PropertyPathToScopePathPointer(p *string) []string {

	parts := strings.Split(*p, "/")
	newParts := []string{"#"}

	for _, part := range parts {
		if part != "" {
			newParts = append(newParts, "properties")
			newParts = append(newParts, part)
		}
	}

	return newParts
}

func filteredPath(scope *string) []string {
	pathParts := strings.Split(*scope, "/")

	filteredParts := []string{}
	for _, part := range pathParts {
		if part != "#" && part != "properties" {
			filteredParts = append(filteredParts, part)
		}
	}
	return filteredParts
}

func FormValuesToGroupedMap(formParams *url.Values, schemaHandler *SchemaHandler) map[string]interface{} {
	formVal := map[string]interface{}{}

	for scope := range *formParams {
		scopePath := filteredPath(&scope)

		var curLayer interface{} = formVal
		for i, p := range scopePath {
			if i == len(scopePath)-1 {
				if l, ok := curLayer.(map[string]interface{}); ok {
					value := formParams.Get(scope)

					// Todo - investigate Not setting map value to trigger "required check" from schema
					// Todo - parse string value from Post to appropriate int/float/bool/date CORRECTLY
					l[p] = value
					// switch (*schemaHandler).GetValueTypeAtScope(&scope) {
					// case "integer":
					// 	{
					// 		i, err := strconv.Atoi(value)
					// 		if err != nil {
					// 			fmt.Println("can't convert to int", err)
					// 			l[p] = value
					// 		}
					// 		l[p] = i
					// 	}
					// case "number":
					// 	{
					// 		f, err := strconv.ParseFloat(value, 64)
					// 		if err != nil {
					// 			fmt.Println("can't convert to number", err)
					// 			l[p] = value
					// 		}
					//
					// 		fmt.Println("converted to number", f)
					// 		l[p] = f
					// 	}
					// default:
					// 	{
					// 		l[p] = value
					// 	}
					// }
				}
			} else {
				if l, ok := curLayer.(map[string]interface{}); ok {
					innerLayer := l[p]
					if innerLayer == nil {
						l[p] = map[string]interface{}{}

						curLayer = l[p]
					}

				}
			}

		}

	}

	return formVal
}
