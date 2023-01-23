package tags

import (
	"strings"
	"unicode"
)

// Get the tag and inline values
// first non empty tag value is retrieved
// only fieldTag is being checked for inline value (fieldTag: "tag,inline")
func GetTag(fieldTag, bsonTag, name string) (string, bool) {
	fieldArr := SplitTrim(fieldTag, ",")
	n := len(fieldArr)

	isInline := n > 1 && strings.ToLower(fieldArr[1]) == "inline"
	if n > 0 && fieldArr[0] != "" {
		return fieldArr[0], isInline
	}

	bsonArr := SplitTrim(bsonTag, ",")
	if len(bsonArr) > 0 && bsonArr[0] != "" {
		return bsonArr[0], isInline
	}

	return firstCharLower(name), isInline
}

// Converts the first character of a string to lower case
func firstCharLower(name string) string {
	a := []rune(name)
	a[0] = unicode.ToLower(a[0])
	return string(a)
}
