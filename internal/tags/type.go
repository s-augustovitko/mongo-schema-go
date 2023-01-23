package tags

import (
	"fmt"
	"reflect"
)

var validBsonTypes = map[string]bool{
	"string":              true,
	"double":              true,
	"object":              true,
	"array":               true,
	"binData":             true,
	"undefined":           true,
	"objectId":            true,
	"bool":                true,
	"date":                true,
	"null":                true,
	"regex":               true,
	"dbPointer":           true,
	"javascript":          true,
	"symbol":              true,
	"javascriptWithScope": true,
	"int":                 true,
	"timestamp":           true,
	"long":                true,
	"decimal":             true,
	"minKey":              true,
	"maxKey":              true,
}

var bsonMap = map[reflect.Kind]string{
	reflect.Int:       "int",
	reflect.Int8:      "int",
	reflect.Int32:     "int",
	reflect.Int64:     "long",
	reflect.Float32:   "double",
	reflect.Float64:   "decimal",
	reflect.String:    "string",
	reflect.Bool:      "bool",
	reflect.Array:     "array",
	reflect.Slice:     "array",
	reflect.Struct:    "object",
	reflect.Map:       "object",
	reflect.Interface: "objectId",
}

// Get type from type tag or reflect.Kind
func GetType(typeTag string, kind reflect.Kind) ([]string, error) {
	typeArr, err := checkValidTypeArr(SplitTrim(typeTag, ","))
	if len(typeArr) > 0 || err != nil {
		return typeArr, err
	}

	objType, ok := bsonMap[kind]
	if !ok {
		return typeArr, fmt.Errorf("type [%v] is not supported", kind)
	}

	return []string{objType}, nil
}

// Checks that all values in the typeArr are valid bson types
func checkValidTypeArr(typeArr []string) ([]string, error) {
	out := []string{}
	invalid := []string{}
	for _, val := range typeArr {
		if val == "" {
			continue
		}

		if found := validBsonTypes[val]; found {
			out = append(out, val)
		} else {
			invalid = append(invalid, val)
		}
	}

	if len(invalid) > 0 {
		return out, fmt.Errorf("the following types are invalid %v", invalid)
	}
	return out, nil
}
