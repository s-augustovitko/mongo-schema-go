package validation

import (
	"fmt"

	"github.com/s-augustovitko/mongo-schema-go/internal/tags"
)

type ErrorWithTag interface {
	Name() string
	Tag() string
	Error() string
}

type errorWithTag struct {
	name  string
	tag   string
	error string
}

// Get Name value
func (e errorWithTag) Name() string {
	return e.name
}

// Get Tag value
// If tag is empty name value will be retrieved with the first character lower cased
func (e errorWithTag) Tag() string {
	if e.tag == "" {
		val, _ := tags.GetTag("", "", e.name)
		return val
	}
	return e.tag
}

// Gets the error message
func (e errorWithTag) Error() string {
	return fmt.Sprintf("[%v]: %v", e.Name(), e.error)
}

// Creates error with tag and name
// tag can be an empty string
func createErrorWithTag(tag, name string, err error) ErrorWithTag {
	return errorWithTag{
		tag:   tag,
		name:  name,
		error: err.Error(),
	}
}
