package tags

import "testing"

type firstCharLowerTest struct {
	arg  string
	want string
}

func TestFirstCharLower(t *testing.T) {
	tests := []firstCharLowerTest{
		{"SomeCoolName", "someCoolName"},
		{"name1", "name1"},
		{"1Name1", "1Name1"},
	}

	for _, test := range tests {
		if have := firstCharLower(test.arg); have != test.want {
			t.Errorf("\nGot: %#v;\nWant: %#v", have, test.want)
		}
	}
}

type getTagTest struct {
	arg1, arg2, arg3 string
	want             string
	wantInline       bool
}

func TestGetTag(t *testing.T) {
	tests := []getTagTest{
		{"name", "", "", "name", false},
		{"", "name", "", "name", false},
		{"", "", "Name", "name", false},
		{"  name   ,  inline  ", "name2", "Name3", "name", true},
		{"  ,  inline  ", "name2", "Name3", "name2", true},
		{"  ,  inline  ", "", "Name3", "name3", true},
		{" name ,    ", "", "Name3", "name", false},
		{"name,", "name2", "Name3", "name", false},
		{"Name123", "name2", "Name3", "Name123", false},
	}

	for _, test := range tests {
		have, haveInline := GetTag(test.arg1, test.arg2, test.arg3)

		if have != test.want || test.wantInline != haveInline {
			t.Errorf("\nGot: %#v, %#v;\nWant: %#v, %#v", have, haveInline, test.want, test.wantInline)
		}
	}
}
