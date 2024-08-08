package cmd

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

var Variables = make(map[string]interface{}) // Exported variable

// VariableCommand handles the variable command and stores the variable in the map
func VariableCommand(attributes string) {
	attrMap := ParseAttributes(attributes)

	name := attrMap["name"]
	value := attrMap["value"]
	variableType := attrMap["type"]
	var parsedValue interface{}

	switch variableType {
	case "string":
		parsedValue = value
	case "int":
		if i, err := strconv.Atoi(value); err == nil {
			parsedValue = i
		} else {
			parsedValue = value
		}
	case "float":
		if f, err := strconv.ParseFloat(value, 64); err == nil {
			parsedValue = f
		} else {
			parsedValue = value
		}
	case "boolean":
		parsedValue = value == "true"
	case "array":
		var arrayValue []interface{}
		err := json.Unmarshal([]byte(strings.Replace(value, "'", `"`, -1)), &arrayValue)
		if err != nil {
			parsedValue = nil
		} else {
			parsedValue = arrayValue
		}
	case "object":
		var objectValue map[string]interface{}
		err := json.Unmarshal([]byte(strings.Replace(value, "'", `"`, -1)), &objectValue)
		if err != nil {
			parsedValue = nil
		} else {
			// Convert number strings to actual numbers
			for k, v := range objectValue {
				if strVal, ok := v.(string); ok {
					if i, err := strconv.Atoi(strVal); err == nil {
						objectValue[k] = i
					} else if f, err := strconv.ParseFloat(strVal, 64); err == nil {
						objectValue[k] = f
					}
				}
			}
			parsedValue = objectValue
		}
	case "date":
		parsedValue = value
	}

	if name != "" {
		Variables[name] = parsedValue
		fmt.Printf("Variable set: %s = %v (type: %s)\n", name, parsedValue, variableType)
	} else {
		fmt.Printf("Warning: Attempted to set variable with empty name. Value: %v, Type: %s\n", value, variableType)
	}
}
