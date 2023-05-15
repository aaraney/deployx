package commands

import (
	"fmt"
	"os"
	"strings"
	"unicode"
)

// validateStackName checks if the provided string is a valid stack name (namespace).
// It currently only does a rudimentary check if the string is empty, or consists
// of only whitespace and quoting characters.
func validateStackName(namespace string) error {
	v := strings.TrimFunc(namespace, quotesOrWhitespace)
	if v == "" {
		return fmt.Errorf("invalid stack name: %q", namespace)
	}
	return nil
}

func quotesOrWhitespace(r rune) bool {
	return unicode.IsSpace(r) || r == '"' || r == '\''
}

// merge the right map into the left map. intersecting keys are returned but not replaced in the
// left map.
func leftMerge(left, right map[string]string) []string {
	right_in_left := []string{}

	for key, value := range right {
		if _, ok := left[key]; ok {
			right_in_left = append(right_in_left, value)
			continue
		}
		left[key] = value
	}

	return right_in_left
}

// retrieves environment variables, split them on '=', and returns them as map[string]string.
func environMap() map[string]string {
	env := os.Environ()
	env_map := make(map[string]string)
	for _, variable := range env {
		split_loc := strings.IndexRune(variable, '=')
		variable_name := variable[:split_loc]
		variable_value := variable[split_loc+1:]
		env_map[variable_name] = variable_value
	}
	return env_map
}
