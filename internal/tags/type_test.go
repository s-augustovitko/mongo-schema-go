package tags

import (
	"reflect"
	"testing"
)

type checkValidTypeArrTest struct {
	arg, want []string
	haveErr   bool
}

func TestCheckValidTypeArr(t *testing.T) {
	tests := []checkValidTypeArrTest{
		{[]string{"string", "decimal"}, []string{"string", "decimal"}, false},
		{[]string{"invalid", "double"}, []string{"double"}, true},
		{[]string{"invalid"}, []string{}, true},
	}

	for _, test := range tests {
		have, err := checkValidTypeArr(test.arg)

		if !CompareArr(have, test.want) || (test.haveErr && err == nil) {
			t.Errorf("\nGot: %#v;\nWant: %#v;\nErr: %#v", have, test.want, err)
		}
	}
}

type getTypeTest struct {
	arg1    string
	arg2    reflect.Kind
	want    []string
	haveErr bool
}

func TestGetType(t *testing.T) {
	tests := []getTypeTest{
		{"", reflect.String, []string{"string"}, false},
		{",  ", reflect.String, []string{"string"}, false},
		{"invalid", reflect.String, []string{}, true},
		{",decimal", reflect.String, []string{"decimal"}, false},
		{"", reflect.Chan, []string{}, true},
		{" bool, double  ", reflect.String, []string{"bool", "double"}, false},
	}

	for _, test := range tests {
		have, err := GetType(test.arg1, test.arg2)

		if !CompareArr(have, test.want) || (test.haveErr && err == nil) {
			t.Errorf("\nGot: %#v;\nWant: %#v;\nErr: %#v", have, test.want, err)
		}
	}
}
