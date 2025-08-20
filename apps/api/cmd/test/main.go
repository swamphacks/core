package main

import (
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"strconv"
)

var (
	ContentKey       = "content"
	FormItemTypeKey  = "type"
	FormFieldNameKey = "name"
)

func main() {
	fileData, err := os.ReadFile("./applicationFormExample.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	var data map[string]any

	err = json.Unmarshal(fileData, &data)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(*buildValidationStruct(data[ContentKey]))

	// https://pkg.go.dev/github.com/go-playground/validator/v10#Validate.Var
}

func buildValidationStruct(content any) *map[string]string {
	contentArray, ok := content.([]any)

	if !ok {
		return nil
	}

	validationStruct := make(map[string]string)

	for _, v := range contentArray {
		item, ok := v.(map[string]any)

		if !ok {
			continue
		}

		if item[FormItemTypeKey] == nil {
			fieldName, ok := item[FormFieldNameKey].(string)
			if !ok {
				continue
			}

			validationStruct[fieldName] = ""

			_, ok = item["required"].(bool)

			if !ok {
				validationStruct[fieldName] = ""
			} else {
				validationStruct[fieldName] = "required"
			}

			validation, ok := item["validation"].(map[string]any)

			if !ok {
				continue
			}

			for validationKey := range validation {
				rule := ""
				val := validation[validationKey]

				switch val2 := val.(type) {
				case string:
					rule = fmt.Sprintf("%s=%s", validationKey, val2)
				case float64:
					rule = fmt.Sprintf("%s=%s", validationKey, strconv.FormatFloat(val2, 'f', 0, 64))
				case int:
					rule = fmt.Sprintf("%s=%d", validationKey, val2)
				default:
					continue
				}

				validationStruct[fieldName] += "," + rule
			}
		} else {
			result := buildValidationStruct(item[ContentKey])

			if result != nil {
				maps.Copy(validationStruct, *result)
			}

		}
	}

	return &validationStruct
}
