package validation

import (
	"reflect"
	"testing"
)

type setValTest struct {
	arg  WithVal[string]
	want bool
}

func TestSetVal(t *testing.T) {
	tests := []setValTest{
		{WithVal[string]{Val: "test", Exists: true}, true},
		{WithVal[string]{Val: "", Exists: true}, true},
		{WithVal[string]{Val: "test", Exists: false}, false},
		{WithVal[string]{Val: "", Exists: false}, false},
	}

	for _, test := range tests {
		obj := BsonM{}
		test.arg.SetVal("field", &obj)

		have, ok := obj["field"]
		if ok != test.want {
			t.Errorf("\nGot: %#v;\nWant: %#v", ok, test.want)
		}
		if ok && have != test.arg.Val {
			t.Errorf("\nGot: %#v;\nWant: %#v", have, test.arg.Val)
		}
	}
}

type addValidationsTest struct {
	arg1 []string
	arg2 Validation
	want BsonM
}

func TestAddValidations(t *testing.T) {
	tests := []addValidationsTest{
		{[]string{}, Validation{},
			BsonM{}},
		{[]string{"string", "int", "long", "bool", "object"}, Validation{Required: true, UniqueItems: true},
			BsonM{}},
		{[]string{"string"}, Validation{Min: CreateVal(1.3), Max: CreateVal(10.2)},
			BsonM{"maxLength": 10, "minLength": 1}},
		{[]string{"string"}, Validation{PatternProps: CreateVal("gi")},
			BsonM{}},
		{[]string{"string"}, Validation{Pattern: CreateVal("@gmail.com$"), PatternProps: CreateVal("gi")},
			BsonM{"pattern": "@gmail.com$", "patternProperties": "gi"}},
		{[]string{"int"}, Validation{Min: CreateVal[float64](1), Max: CreateVal[float64](20), MultipleOf: CreateVal[float64](4)},
			BsonM{"maximum": float64(20), "minimum": float64(1), "multipleOf": float64(4)}},
		{[]string{"long"}, Validation{Min: CreateVal(1.1), Max: CreateVal(20.4), MultipleOf: CreateVal(4.2)},
			BsonM{"maximum": 20.4, "minimum": 1.1, "multipleOf": 4.2}},
		{[]string{"double"}, Validation{Min: CreateVal(1.1), Max: CreateVal(20.4), MultipleOf: CreateVal(4.2)},
			BsonM{"maximum": 20.4, "minimum": 1.1, "multipleOf": 4.2}},
		{[]string{"decimal"}, Validation{Min: CreateVal(1.1), Max: CreateVal(20.4), MultipleOf: CreateVal(4.2)},
			BsonM{"maximum": 20.4, "minimum": 1.1, "multipleOf": 4.2}},
		{[]string{"array"}, Validation{Min: CreateVal(1.1), Max: CreateVal(20.4), UniqueItems: true},
			BsonM{"maxItems": 20, "minItems": 1, "uniqueItems": true}},
		{[]string{"object"}, Validation{Min: CreateVal(1.1), Max: CreateVal(20.4)},
			BsonM{"maxProperties": 20, "minProperties": 1}},
	}

	for _, test := range tests {
		have := BsonM{}
		addValidations(test.arg1, test.arg2, &have)

		if !reflect.DeepEqual(have, test.want) {
			t.Errorf("\nGot: %#v;\nWant: %#v;", have, test.want)
		}
	}
}

type parseValidationTest struct {
	arg     string
	want    Validation
	wantErr bool
}

func TestParseValidation(t *testing.T) {
	tests := []parseValidationTest{
		{"", Validation{}, false},
		{" ,   ", Validation{}, false},
		{"invalid=false", Validation{}, true},
		{"required", Validation{Required: true}, false},
		{"required=4", Validation{}, true},
		{"uniqueItems", Validation{UniqueItems: true}, false},
		{"uniqueItems=4", Validation{}, true},
		{"min", Validation{}, true},
		{"min=4.4", Validation{Min: CreateVal(4.4)}, false},
		{"min=asd", Validation{}, true},
		{"max", Validation{}, true},
		{"max=6.3", Validation{Max: CreateVal(6.3)}, false},
		{"max=5.0,min=6.1", Validation{Max: CreateVal(5.0), Min: CreateVal(6.1)}, true},
		{"max=5.0,min=5.0", Validation{Max: CreateVal(5.0), Min: CreateVal(5.0)}, false},
		{"max=asd", Validation{}, true},
		{"multipleOf", Validation{}, true},
		{"multipleOf=5", Validation{MultipleOf: CreateVal[float64](5)}, false},
		{"multipleOf=asd", Validation{}, true},
		{"pattern", Validation{}, true},
		{"pattern=@gmail.com$", Validation{Pattern: CreateVal("@gmail.com$")}, false},
		{"patternProperties", Validation{}, true},
		{"patternProperties=gi", Validation{PatternProps: CreateVal("gi")}, false},
	}

	for _, test := range tests {
		have, err := parseValidation(test.arg)
		haveErr := err != nil
		if !reflect.DeepEqual(have, test.want) || test.wantErr != haveErr {
			t.Errorf("\nGot: %#v;\nWant: %#v;\nErr: %#v", have, test.want, err)
		}
	}
}
