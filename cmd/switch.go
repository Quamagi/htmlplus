package cmd

// SwitchCommand handles the switch-case logic
func SwitchCommand(variable string, cases map[string]string, defaultContent string) string {
	value := GetVariable(variable)

	if caseContent, exists := cases[value]; exists {
		return caseContent
	}

	return defaultContent
}
