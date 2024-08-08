package cmd

import (
	"fmt"
	"strings"
)

// PrintCommand handles the print command and wraps output in the specified HTML tag and classes
func PrintCommand(attributes string) string {
	attrMap := ParseAttributes(attributes)

	var valor string
	if name, exists := attrMap["name"]; exists {
		if property, propExists := attrMap["property"]; propExists {
			// Obtener valor de la propiedad anidada
			parts := strings.Split(name, ".")
			parts = append(parts, property)
			valor = toString(GetNestedVariable(parts))
		} else {
			parts := strings.Split(name, ".")
			valor = toString(GetNestedVariable(parts))
		}
		fmt.Printf("PrintCommand: Retrieved value for %s = %s\n", name, valor)
	} else {
		valor = attrMap["value"]
		fmt.Printf("PrintCommand: Using literal value: %s\n", valor)
	}

	tag := "span" // Default tag
	if customTag, exists := attrMap["tag"]; exists {
		tag = customTag
	}

	classAttr := ""
	if class, exists := attrMap["class"]; exists {
		classAttr = fmt.Sprintf(` class="%s"`, class)
	}

	idAttr := ""
	if id, exists := attrMap["id"]; exists {
		idAttr = fmt.Sprintf(` id="%s"`, id)
	}

	return fmt.Sprintf("<%s%s%s>%s</%s>", tag, classAttr, idAttr, valor, tag)
}
