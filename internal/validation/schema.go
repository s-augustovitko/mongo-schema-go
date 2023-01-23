package validation

import (
	"fmt"
	"reflect"

	"github.com/s-augustovitko/mongo-schema-go/internal/tags"
)

const (
	tagBson      = "bson"
	tagField     = "field"
	tagType      = "type"
	tagItemsType = "itemsType"
	tagValid     = "validation"
	tagDesc      = "description"
	tagEnum      = "enum"
	tagItems     = "items"
)

type config struct {
	Validation      Validation
	Tag             string
	BsonType        []string
	Enum            WithVal[[]string]
	Description     WithVal[string]
	ItemsBsonType   []string
	ItemsValidation Validation
	IsArray         bool
	IsArrayOfStruct bool
	IsStruct        bool
	IsInline        bool
}

// Creates the configuration used for the json schema parsing
// config is created based on the reflect value of a field
func createConfig(value reflect.Value, field reflect.StructField) (config, error) {
	var err error
	cfg := config{}

	// FIELD
	cfg.IsStruct = value.Kind() == reflect.Struct
	cfg.Tag, cfg.IsInline = tags.GetTag(field.Tag.Get(tagField), field.Tag.Get(tagBson), field.Name)
	cfg.IsInline = cfg.IsStruct && cfg.IsInline // IsInline can only be true if the field is a struct
	// DESCRIPTION
	description := field.Tag.Get(tagDesc)
	if description != "" {
		cfg.Description = CreateVal(description)
	}
	// ENUM
	enum := tags.SplitTrim(field.Tag.Get(tagEnum), ",")
	if len(enum) > 0 {
		cfg.Enum = CreateVal(enum)
	}
	// TYPE
	cfg.BsonType, err = tags.GetType(field.Tag.Get(tagType), value.Kind())
	if err != nil {
		return cfg, err
	}
	// VALIDATION
	cfg.Validation, err = parseValidation(field.Tag.Get(tagValid))
	cfg.Validation.Required = !cfg.IsInline && cfg.Validation.Required // Required can not be set if it is inline
	if err != nil {
		return cfg, err
	}

	// STRUCTS AND ARRAYS
	cfg.IsArray = value.Kind() == reflect.Array || value.Kind() == reflect.Slice
	if !cfg.IsArray {
		return cfg, nil
	}

	// ARRAYS
	if value.Len() < 1 {
		return cfg, fmt.Errorf("could not properly parse the slice since its empty")
	}
	item := value.Index(0)
	if item.Kind() == reflect.Pointer {
		item = item.Elem()
	}

	// ARRAY OF STRUCTS
	if item.Kind() == reflect.Struct {
		cfg.IsArrayOfStruct = true
		return cfg, nil
	}

	// ARRAYS
	cfg.ItemsBsonType, err = tags.GetType(field.Tag.Get(tagItemsType), item.Kind())
	if err != nil {
		return cfg, err
	}
	cfg.ItemsValidation, err = parseValidation(field.Tag.Get(tagItems))
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}

// Creates the json schema from the reflect of the struct
// Struct can be empty except for arrays, all arrays must be filled with at least 1 item of its kind
// Returns an array of required fields.([]string) and warnings.(ErrorWithTag)
func CreateJSONSchema(value reflect.Value, objProperties *BsonM) ([]string, []error) {
	requiredFields := []string{}
	errors := []error{}

	valTyp := value.Type()
	for i := 0; i < value.NumField(); i++ {
		val := value.Field(i)
		field := valTyp.Field(i)
		if val.Kind() == reflect.Pointer {
			val = val.Elem()
		}

		// CONFIG
		cfg, err := createConfig(val, field)
		if err != nil {
			errors = append(errors, createErrorWithTag(cfg.Tag, field.Name, err))
			continue
		}

		// BASE VALUES
		obj := BsonM{"bsonType": cfg.BsonType}
		addValidations(cfg.BsonType, cfg.Validation, &obj)
		if cfg.Validation.Required {
			requiredFields = append(requiredFields, cfg.Tag)
		}

		// INLINE STRUCT
		if cfg.IsStruct && cfg.IsInline {
			props := BsonM{}
			reqs, errs := CreateJSONSchema(val, &props)

			for k, v := range props {
				(*objProperties)[k] = v
			}
			requiredFields = append(requiredFields, reqs...)
			errors = append(errors, errs...)
			continue
		}

		// STRUCT
		if cfg.IsStruct {
			props := BsonM{}
			reqs, errs := CreateJSONSchema(val, &props)
			obj["properties"] = props
			obj["required"] = reqs
			errors = append(errors, errs...)
			(*objProperties)[cfg.Tag] = obj
			continue
		}

		// ARRAY OF STRUCTS
		if cfg.IsArrayOfStruct {
			props := BsonM{}
			reqs, errs := CreateJSONSchema(val.Index(0), &props)
			obj["items"] = BsonM{
				"bsonType":   []string{"object"},
				"required":   reqs,
				"properties": props,
			}
			errors = append(errors, errs...)
			(*objProperties)[cfg.Tag] = obj
			continue
		}

		// ARRAY
		if cfg.IsArray {
			items := BsonM{"bsonType": cfg.ItemsBsonType}
			cfg.Enum.SetVal("enum", &items)
			addValidations(cfg.ItemsBsonType, cfg.ItemsValidation, &items)
			cfg.Description.SetVal("description", &items)

			obj["items"] = items
		} else {
			cfg.Description.SetVal("description", &obj)
			cfg.Enum.SetVal("enum", &obj)
		}

		(*objProperties)[cfg.Tag] = obj
	}

	return requiredFields, errors
}
