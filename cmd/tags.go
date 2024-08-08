package cmd

import (
	"regexp"
	"strings"
)

// ProcessCustomTags processes the custom HTML tags and executes the corresponding commands
func ProcessCustomTags(content string) string {
	// Process <variable> tags
	reVariable := regexp.MustCompile(`(?s)<variable\s+(.*?)>(.*?)<\/variable>`)
	content = reVariable.ReplaceAllStringFunc(content, func(match string) string {
		attributes := reVariable.FindStringSubmatch(match)[1]
		VariableCommand(attributes)
		return ""
	})

	// Process <function> tags
	reFunction := regexp.MustCompile(`(?s)<function\s+(.*?)>(.*?)<\/function>`)
	content = reFunction.ReplaceAllStringFunc(content, func(match string) string {
		attributes := reFunction.FindStringSubmatch(match)[1]
		innerContent := reFunction.FindStringSubmatch(match)[2]
		DefineFunctionCommand(attributes, innerContent)
		return ""
	})

	// Process <call-function> tags
	reCallFunction := regexp.MustCompile(`(?s)<call-function\s+(.*?)>(.*?)<\/call-function>`)
	content = reCallFunction.ReplaceAllStringFunc(content, func(match string) string {
		attributes := reCallFunction.FindStringSubmatch(match)[1]
		innerContent := reCallFunction.FindStringSubmatch(match)[2]
		return CallFunctionCommand(attributes, innerContent)
	})

	// Process <operation> tags
	reOperation := regexp.MustCompile(`(?s)<operation\s+(.*?)>(.*?)<\/operation>`)
	content = reOperation.ReplaceAllStringFunc(content, func(match string) string {
		attributes := reOperation.FindStringSubmatch(match)[1]
		operandsMatches := regexp.MustCompile(`(?s)<operand>(.*?)<\/operand>`).FindAllStringSubmatch(reOperation.FindStringSubmatch(match)[2], -1)
		var operands []string
		for _, operandMatch := range operandsMatches {
			operands = append(operands, strings.TrimSpace(operandMatch[1]))
		}
		OperationCommand(attributes, operands)
		return ""
	})

	// Process <if> tags
	reIf := regexp.MustCompile(`(?s)<if\s+(.*?)>(.*?)<\/if>`)
	content = reIf.ReplaceAllStringFunc(content, func(match string) string {
		attributes := reIf.FindStringSubmatch(match)[1]
		innerContent := reIf.FindStringSubmatch(match)[2]

		thenMatch := regexp.MustCompile(`(?s)<then>(.*?)<\/then>`).FindStringSubmatch(innerContent)
		elseMatch := regexp.MustCompile(`(?s)<else>(.*?)<\/else>`).FindStringSubmatch(innerContent)

		var thenContent, elseContent string
		if len(thenMatch) > 1 {
			thenContent = thenMatch[1]
		}
		if len(elseMatch) > 1 {
			elseContent = elseMatch[1]
		}

		return ConditionCommand(attributes, thenContent, elseContent)
	})

	// Process <switch> tags
	reSwitch := regexp.MustCompile(`(?s)<switch\s+(.*?)>(.*?)<\/switch>`)
	content = reSwitch.ReplaceAllStringFunc(content, func(match string) string {
		attributes := reSwitch.FindStringSubmatch(match)[1]
		innerContent := reSwitch.FindStringSubmatch(match)[2]

		variable := ParseAttributes(attributes)["variable"]

		reCase := regexp.MustCompile(`(?s)<case\s+value="(.*?)">(.*?)<\/case>`)
		cases := make(map[string]string)
		caseMatches := reCase.FindAllStringSubmatch(innerContent, -1)
		for _, caseMatch := range caseMatches {
			cases[caseMatch[1]] = caseMatch[2]
		}

		defaultContent := ""
		reDefault := regexp.MustCompile(`(?s)<default>(.*?)<\/default>`)
		defaultMatch := reDefault.FindStringSubmatch(innerContent)
		if len(defaultMatch) > 1 {
			defaultContent = defaultMatch[1]
		}

		return SwitchCommand(variable, cases, defaultContent)
	})

	// Process <for-each> tags
	reForEach := regexp.MustCompile(`(?s)<for-each\s+(.*?)>(.*?)<\/for-each>`)
	content = reForEach.ReplaceAllStringFunc(content, func(match string) string {
		attributes := reForEach.FindStringSubmatch(match)[1]
		innerContent := reForEach.FindStringSubmatch(match)[2]
		return ForEachCommand(attributes, innerContent)
	})

	// Process <fetch> tags
	reFetch := regexp.MustCompile(`(?s)<fetch\s+(.*?)>(.*?)<\/fetch>`)
	content = reFetch.ReplaceAllStringFunc(content, func(match string) string {
		attributes := reFetch.FindStringSubmatch(match)[1]
		innerContent := reFetch.FindStringSubmatch(match)[2]

		successMatch := regexp.MustCompile(`(?s)<on-success>(.*?)<\/on-success>`).FindStringSubmatch(innerContent)
		errorMatch := regexp.MustCompile(`(?s)<on-error>(.*?)<\/on-error>`).FindStringSubmatch(innerContent)

		var successContent, errorContent string
		if len(successMatch) > 1 {
			successContent = successMatch[1]
		}
		if len(errorMatch) > 1 {
			errorContent = errorMatch[1]
		}

		return FetchCommand(attributes, successContent, errorContent)
	})

	// Process <print> tags
	rePrint := regexp.MustCompile(`(?s)<print\s+(.*?)>(.*?)<\/print>`)
	content = rePrint.ReplaceAllStringFunc(content, func(match string) string {
		attributes := rePrint.FindStringSubmatch(match)[1]
		return PrintCommand(attributes)
	})

	return content
}
