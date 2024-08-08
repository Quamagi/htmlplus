package cmd

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// ParseAttributes parses the attributes of a custom HTML tag into a map
func ParseAttributes(attributeString string) map[string]string {
	attrMap := make(map[string]string)
	attributes := strings.Split(attributeString, `"`)

	for i := 0; i < len(attributes)-1; i += 2 {
		key := strings.TrimSpace(strings.Trim(attributes[i], " ="))
		value := attributes[i+1]
		attrMap[key] = value
	}

	return attrMap
}

// GetVariable retrieves the value of a variable by name
func GetVariable(name string) string {
	if value, exists := Variables[name]; exists {
		return fmt.Sprintf("%v", value)
	}
	return "null"
}

// GetNestedVariable retrieves a nested variable by its path
func GetNestedVariable(path []string) interface{} {
	var value interface{} = Variables
	for _, p := range path {
		switch v := value.(type) {
		case map[string]interface{}:
			value = v[p]
		case []interface{}:
			index, err := strconv.Atoi(p)
			if err != nil {
				return nil
			}
			value = v[index]
		default:
			return nil
		}
	}
	return value
}

// toString converts an interface to a string
func toString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case json.Number:
		return v.String()
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case int:
		return strconv.Itoa(v)
	case bool:
		if v {
			return "true"
		}
		return "false"
	case []interface{}:
		var strValues []string
		for _, item := range v {
			strValues = append(strValues, toString(item))
		}
		return "[" + strings.Join(strValues, ", ") + "]"
	case map[string]interface{}:
		var strPairs []string
		for key, val := range v {
			strPairs = append(strPairs, fmt.Sprintf("%s: %s", key, toString(val)))
		}
		return "{" + strings.Join(strPairs, ", ") + "}"
	default:
		return fmt.Sprintf("%v", v)
	}
}
