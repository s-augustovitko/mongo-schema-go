package tags

import (
	"testing"
)

type splitTrimTest struct {
	arg1 string
	arg2 string
	want []string
}

func TestSplitTrim(t *testing.T) {
	tests := []splitTrimTest{
		{"", ",", []string{}},
		{"asd", ",", []string{"asd"}},
		{"asd, asd, asd", ",", []string{"asd", "asd", "asd"}},
		{"  , asd, asd, ", ",", []string{"", "asd", "asd", ""}},
		{"asd= asd  =", ",", []string{"asd= asd  ="}},
		{"asd= asd  =", "=", []string{"asd", "asd", ""}},
	}

	for _, test := range tests {
		have := SplitTrim(test.arg1, test.arg2)

		if !CompareArr(have, test.want) {
			t.Errorf("\nGot: %#v;\nWant: %#v", have, test.want)
		}
	}
}

type compareArrTest struct {
	arg1, arg2 []string
	want       bool
}

func TestCompareArr(t *testing.T) {
	tests := []compareArrTest{
		{[]string{}, []string{}, true},
		{[]string{"a", "b", "c"}, []string{"a", "b", "c"}, true},
		{[]string{"A", "b", "1"}, []string{"A", "B", "1"}, false},
		{[]string{"a", "c", "b"}, []string{"a", "b", "c"}, false},
		{[]string{"a"}, []string{"a", "b"}, false},
	}

	for _, test := range tests {
		if have := CompareArr(test.arg1, test.arg2); have != test.want {
			t.Errorf("\nGot: %#v;\nWant: %#v;\nTest: %#v", have, test.want, test)
		}
	}
}
