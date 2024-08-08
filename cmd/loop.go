package cmd

import (
	"fmt"
	"strings"
)

// ForEachCommand maneja el bucle for-each
func ForEachCommand(attributes string, innerContent string) string {
	attrMap := ParseAttributes(attributes)
	arrayName := attrMap["array"]
	itemName := attrMap["item"]

	parts := strings.Split(arrayName, ".")
	array := GetNestedVariable(parts)

	var result string
	switch items := array.(type) {
	case []interface{}:
		for i, item := range items {
			Variables[itemName] = item
			fmt.Printf("ForEachCommand: Processing item %d in array %s\n", i, arrayName)
			result += ProcessCustomTags(innerContent)
			delete(Variables, itemName)
		}
	default:
		fmt.Printf("ForEachCommand: %s is not an array\n", arrayName)
	}

	return result
}
