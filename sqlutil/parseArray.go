package sqlutil

import (
	"fmt"
	"regexp"
	"strings"
)

// https://gist.github.com/adharris/4163702

// PARSING ARRAYS
// SEE http://www.postgresql.org/docs/9.1/static/arrays.html#ARRAYS-IO
// Arrays are output within {} and a delimiter, which is a comma for most
// postgres types (; for box)
//
// Individual values are surrounded by quotes:
// The array output routine will put double quotes around element values if
// they are empty strings, contain curly braces, delimiter characters,
// double quotes, backslashes, or white space, or match the word NULL.
// Double quotes and backslashes embedded in element values will be
// backslash-escaped. For numeric data types it is safe to assume that double
// quotes will never appear, but for textual data types one should be prepared
// to cope with either the presence or absence of quotes.

// construct a regexp to extract values:
var (
	// unquoted array values must not contain: (" , \ { } whitespace NULL)
	// and must be at least one char
	unquotedChar  = `[^",\\{}\s(NULL)]`
	unquotedValue = fmt.Sprintf("(%s)+", unquotedChar)

	// quoted array values are surrounded by double quotes, can be any
	// character except " or \, which must be backslash escaped:
	quotedChar  = `[^"\\]|\\"|\\\\`
	quotedValue = fmt.Sprintf("\"(%s)*\"", quotedChar)

	// an array value may be either quoted or unquoted:
	arrayValue = fmt.Sprintf("(?P<value>(%s|%s))", unquotedValue, quotedValue)

	// Array values are separated with a comma IF there is more than one value:
	arrayExp = regexp.MustCompile(fmt.Sprintf("((%s)(,)?)", arrayValue))

	valueIndex int
)

// Find the index of the 'value' named expression
func init() {
	for i, subexp := range arrayExp.SubexpNames() {
		if subexp == "value" {
			valueIndex = i
			break
		}
	}
}

// Parse the output string from the array type.
// Regex used: (((?P<value>(([^",\\{}\s(NULL)])+|"([^"\\]|\\"|\\\\)*")))(,)?)
func parseArray(array string) []string {
	results := make([]string, 0)
	matches := arrayExp.FindAllStringSubmatch(array, -1)
	for _, match := range matches {
		s := match[valueIndex]
		// the string _might_ be wrapped in quotes, so trim them:
		s = strings.Trim(s, "\"")
		results = append(results, s)
	}
	return results
}
