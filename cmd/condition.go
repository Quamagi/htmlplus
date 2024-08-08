package cmd

// ConditionCommand handles the if-then-else logic
func ConditionCommand(attributes string, thenContent string, elseContent string) string {
	attrMap := ParseAttributes(attributes)
	condition := attrMap["condition"]

	conditionValue := GetVariable(condition)

	if conditionValue == "true" {
		return thenContent
	} else {
		return elseContent
	}
}
