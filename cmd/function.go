package cmd

import (
	"fmt"
	"regexp"
)

// Function representa una función personalizada definida en el HTML
type Function struct {
	Params map[string]string
	Body   string
}

// Almacena todas las funciones definidas
var functions = make(map[string]Function)

// DefineFunctionCommand maneja la definición de una función personalizada
func DefineFunctionCommand(attributes string, innerContent string) {
	attrMap := ParseAttributes(attributes)
	functionName := attrMap["name"]

	paramMap := make(map[string]string)
	reParam := regexp.MustCompile(`(?s)<param\s+(.*?)></param>`)
	innerContent = reParam.ReplaceAllStringFunc(innerContent, func(match string) string {
		paramAttributes := reParam.FindStringSubmatch(match)[1]
		paramAttrMap := ParseAttributes(paramAttributes)
		paramName := paramAttrMap["name"]
		paramType := paramAttrMap["type"]
		paramMap[paramName] = paramType
		return ""
	})

	reBody := regexp.MustCompile(`(?s)<body>(.*?)</body>`)
	bodyContent := reBody.FindStringSubmatch(innerContent)[1]

	functions[functionName] = Function{Params: paramMap, Body: bodyContent}
	fmt.Printf("Function defined: %s with params %v\n", functionName, paramMap)
}

// CallFunctionCommand maneja la llamada a una función personalizada
func CallFunctionCommand(attributes string, innerContent string) string {
	attrMap := ParseAttributes(attributes)
	functionName := attrMap["name"]

	function, exists := functions[functionName]
	if !exists {
		return fmt.Sprintf("<!-- Error: function %s not found -->", functionName)
	}

	argMap := make(map[string]interface{})
	reArg := regexp.MustCompile(`(?s)<arg\s+(.*?)></arg>`)
	for _, match := range reArg.FindAllStringSubmatch(innerContent, -1) {
		argAttributes := match[1]
		argAttrMap := ParseAttributes(argAttributes)
		argName := argAttrMap["name"]
		argValue := argAttrMap["value"]
		argMap[argName] = argValue
	}

	// Guardar las variables originales antes de sobrescribirlas
	originalVariables := make(map[string]interface{})
	for paramName := range function.Params {
		if originalValue, exists := Variables[paramName]; exists {
			originalVariables[paramName] = originalValue
		}
	}

	// Almacenar temporalmente los argumentos en variables para el procesamiento
	for paramName, paramType := range function.Params {
		if argValue, exists := argMap[paramName]; exists {
			VariableCommand(fmt.Sprintf(`name="%s" type="%s" value="%s"`, paramName, paramType, argValue))
			fmt.Printf("Argument set: %s = %s (type: %s)\n", paramName, argValue, paramType)
		} else {
			fmt.Printf("Warning: Argument %s not provided for function %s\n", paramName, functionName)
		}
	}

	// Print all variables before processing function body
	fmt.Println("Variables before processing function body:")
	for name, value := range Variables {
		fmt.Printf("%s: %v\n", name, value)
	}

	result := ProcessCustomTags(function.Body)

	// Restaurar las variables originales
	for paramName := range function.Params {
		if originalValue, exists := originalVariables[paramName]; exists {
			Variables[paramName] = originalValue
		} else {
			delete(Variables, paramName)
		}
	}

	return result
}
