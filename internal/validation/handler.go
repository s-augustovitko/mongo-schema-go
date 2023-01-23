package validation

import (
	"fmt"
	"strconv"

	"github.com/s-augustovitko/mongo-schema-go/internal/tags"
)

type BsonM = map[string]interface{}

type WithVal[T any] struct {
	Exists bool
	Val    T
}

// If item exists then sets the obj[field] to the value
// if (item.Exists) obj[field] = item.Val
func (item WithVal[T]) SetVal(field string, obj *BsonM) {
	if item.Exists {
		(*obj)[field] = item.Val
	}
}

func CreateVal[T any](val T) WithVal[T] {
	return WithVal[T]{Val: val, Exists: true}
}

type Validation struct {
	UniqueItems  bool
	Required     bool
	Min          WithVal[float64]
	Max          WithVal[float64]
	Pattern      WithVal[string]
	PatternProps WithVal[string]
	MultipleOf   WithVal[float64]
}

func parseValidation(validation string) (Validation, error) {
	var out Validation

	validItems := tags.SplitTrim(validation, ",")
	if len(validItems) == 0 {
		return out, nil
	}

	for _, item := range validItems {
		val := tags.SplitTrim(item, "=")
		if len(val) == 0 || val[0] == "" {
			continue
		}

		if val[0] == "required" || val[0] == "uniqueItems" {
			if len(val) != 1 {
				return out, fmt.Errorf("%v validation does not need a value", val)
			}
		} else if len(val) != 2 {
			return out, fmt.Errorf("%v validation requires a value", val)
		}

		switch val[0] {
		case "required":
			out.Required = true
		case "uniqueItems":
			out.UniqueItems = true
		case "min":
			floatVal, err := strconv.ParseFloat(val[1], 64)
			if err != nil {
				return out, err
			}
			out.Min = CreateVal(floatVal)
		case "max":
			floatVal, err := strconv.ParseFloat(val[1], 64)
			if err != nil {
				return out, err
			}
			out.Max = CreateVal(floatVal)
		case "multipleOf":
			floatVal, err := strconv.ParseFloat(val[1], 64)
			if err != nil {
				return out, err
			}
			out.MultipleOf = CreateVal(floatVal)
		case "pattern":
			out.Pattern = CreateVal(val[1])
		case "patternProperties":
			out.PatternProps = CreateVal(val[1])
		default:
			return out, fmt.Errorf("invalid validation value: %v", val)
		}
	}

	if out.Min.Exists && out.Max.Exists && out.Max.Val < out.Min.Val {
		return out, fmt.Errorf("invalid [min,max] values, min can not be larger than max")
	}
	return out, nil
}

func addValidations(types []string, validation Validation, obj *BsonM) {
	for _, kind := range types {
		switch kind {
		case "double", "int", "long", "decimal":
			validation.Max.SetVal("maximum", obj)
			validation.Min.SetVal("minimum", obj)
			validation.MultipleOf.SetVal("multipleOf", obj)
		case "string":
			floatToIntVal(validation.Max).SetVal("maxLength", obj)
			floatToIntVal(validation.Min).SetVal("minLength", obj)
			validation.Pattern.SetVal("pattern", obj)
			if validation.Pattern.Exists {
				validation.PatternProps.SetVal("patternProperties", obj)
			}
		case "array":
			floatToIntVal(validation.Max).SetVal("maxItems", obj)
			floatToIntVal(validation.Min).SetVal("minItems", obj)
			(*obj)["uniqueItems"] = validation.UniqueItems
		case "object":
			floatToIntVal(validation.Max).SetVal("maxProperties", obj)
			floatToIntVal(validation.Min).SetVal("minProperties", obj)
		}
	}
}

// Converts WithVal[float64] to WithVal[int]
func floatToIntVal(item WithVal[float64]) WithVal[int] {
	return WithVal[int]{Val: int(item.Val), Exists: item.Exists}
}
