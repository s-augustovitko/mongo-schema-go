package schema

import (
	"fmt"
	"reflect"

	"github.com/s-augustovitko/mongo-schema-go/internal/validation"
)

// Builds the jsonSchema from a struct
// Struct can be empty except for arrays, all arrays must be filled with at least 1 item of its kind
// Returns: Schema, Warnings (ErrorWithTag), Error
// Any warnings are fields that could not be processed, so they will not show up in the final schema
func Marshal(schema interface{}, title string, additionalProps bool) (out validation.BsonM, warnings []error, err error) {
	if title == "" {
		title = "Schema Validation"
	}

	jsonSchema := validation.BsonM{
		"bsonType":             "object",
		"title":                title,
		"additionalProperties": additionalProps,
	}

	value := reflect.ValueOf(schema)
	if value.Kind() == reflect.Pointer {
		value = value.Elem()
	}
	if value.Kind() != reflect.Struct {
		return jsonSchema, []error{}, fmt.Errorf("to create a validation you must send a struct")
	}

	props := validation.BsonM{}
	reqs, errs := validation.CreateJSONSchema(value, &props)
	jsonSchema["properties"] = props
	jsonSchema["required"] = reqs

	return validation.BsonM{"validator": validation.BsonM{"$jsonSchema": jsonSchema}}, errs, nil
}
