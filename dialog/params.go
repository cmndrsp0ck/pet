package dialog

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/awesome-gocui/gocui"
)

var (
	views = []string{}

	//CurrentCommand is the command before assigning to variables
	CurrentCommand string
	//FinalCommand is the command after assigning to variables
	FinalCommand string

	// This matches parameter patterns like <param> or <param=default>
	// Skips match if there is a whitespace at the end ex. <param='my >
	// Ignores <, > characters since they're used to match the pattern
	// Note: Escaped parameters (\<param\>) are handled separately and excluded from matching
	parameterStringRegex = `<([^<>]*[^\s])>`
)

// insertParams replaces parameter placeholders with actual values while preserving escaped parameters
// Escaped parameters use the syntax \<param\> and are treated as literal text (converted to <param>)
// Regular parameters use the syntax <param> or <param=default> and are replaced with user input
func insertParams(command string, filledInParams map[string]string) string {
	// First, find all escaped patterns (\<param\>) and temporarily replace them with placeholders
	escapedRegex := regexp.MustCompile(`\\<([^>]*)\\>`)
	placeholders := make(map[string]string)
	placeholderCount := 0

	processedCommand := escapedRegex.ReplaceAllStringFunc(command, func(match string) string {
		placeholder := fmt.Sprintf("__ESCAPED_PARAM_%d__", placeholderCount)
		// Extract the inner content and wrap it in <> brackets (removes escape characters)
		submatch := escapedRegex.FindStringSubmatch(match)
		if len(submatch) > 1 {
			content := "<" + submatch[1] + ">"
			placeholders[placeholder] = content
		}
		placeholderCount++
		return placeholder
	})

	// Now find and replace parameters in the processed command
	r := regexp.MustCompile(parameterStringRegex)
	matches := r.FindAllStringSubmatch(processedCommand, -1)

	resultCommand := processedCommand

	for _, p := range matches {
		whole, matchedGroup := p[0], p[1]
		param, _, _ := strings.Cut(matchedGroup, "=")

		// Replace the whole match with the filled-in value of the param
		resultCommand = strings.Replace(resultCommand, whole, filledInParams[param], -1)
	}

	// Restore escaped parameters (now without escape characters)
	for placeholder, original := range placeholders {
		resultCommand = strings.Replace(resultCommand, placeholder, original, -1)
	}

	return resultCommand
}

// removeEscapeChars removes escape characters from escaped parameter patterns
func removeEscapeChars(command string) string {
	// Remove \< and \> escape sequences
	result := strings.ReplaceAll(command, "\\<", "<")
	result = strings.ReplaceAll(result, "\\>", ">")
	return result
}

// SearchForParams returns variables from a command, excluding escaped parameters
// Regular parameters use the syntax <param> or <param=default> and will be extracted for user input
// Escaped parameters use the syntax \<param\> and are treated as literal text, excluded from extraction
// Returns an array of [paramName, defaultValue] pairs
func SearchForParams(command string) [][2]string {
	// First, find all escaped patterns (\<param\>) and temporarily replace them with placeholders
	escapedRegex := regexp.MustCompile(`\\<([^>]*)\\>`)
	placeholders := make(map[string]string)
	placeholderCount := 0

	processedCommand := escapedRegex.ReplaceAllStringFunc(command, func(match string) string {
		placeholder := fmt.Sprintf("__ESCAPED_PARAM_%d__", placeholderCount)
		placeholders[placeholder] = match
		placeholderCount++
		return placeholder
	})

	// Now find parameters in the processed command
	r := regexp.MustCompile(parameterStringRegex)
	matches := r.FindAllStringSubmatch(processedCommand, -1)

	if len(matches) == 0 {
		return nil
	}

	extracted := map[string]string{}
	ordered_params := [][2]string{}

	for _, p := range matches {
		_, matchedGroup := p[0], p[1]

		paramKey, defaultValue, separatorFound := strings.Cut(matchedGroup, "=")
		_, param_exists := extracted[paramKey]

		// Set to empty if no value is provided and param is not already set
		if !separatorFound && !param_exists {
			extracted[paramKey] = ""
		} else if separatorFound {
			// Set to default value instead if it is provided
			extracted[paramKey] = defaultValue
		}

		// Fill in the keys only if seen for the first time to track order
		if !param_exists {
			ordered_params = append(ordered_params, [2]string{paramKey, ""})
		}
	}

	// Fill in the values
	for i, param := range ordered_params {
		pair := [2]string{param[0], extracted[param[0]]}
		ordered_params[i] = pair
	}
	return ordered_params
}

func evaluateParams(g *gocui.Gui, _ *gocui.View) error {
	paramsFilled := map[string]string{}
	for _, v := range views {
		view, _ := g.View(v)
		res := view.Buffer()
		res = strings.Replace(res, "\n", "", -1)
		paramsFilled[v] = strings.TrimSpace(res)
	}
	FinalCommand = insertParams(CurrentCommand, paramsFilled)
	return gocui.ErrQuit
}
