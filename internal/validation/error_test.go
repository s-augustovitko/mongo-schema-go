package validation

import (
	"fmt"
	"testing"
)

func TestCreateError(t *testing.T) {
	want := "[Tag]: Some Error"
	wantTag := "tag"
	wantName := "Tag"

	have := createErrorWithTag(wantTag, wantName, fmt.Errorf("Some Error"))
	if have.Error() != want {
		t.Errorf("\nGot: %#v;\nWant: %#v", have.Error(), want)
	}
	if have.Tag() != wantTag {
		t.Errorf("\nGot: %#v;\nWant: %#v", have.Tag(), wantTag)
	}
	if have.Name() != wantName {
		t.Errorf("\nGot: %#v;\nWant: %#v", have.Name(), wantName)
	}

	have = createErrorWithTag("", wantName, fmt.Errorf("Some Error"))
	if have.Error() != want {
		t.Errorf("\nGot: %#v;\nWant: %#v", have.Error(), want)
	}
	if have.Tag() != wantTag {
		t.Errorf("\nGot: %#v;\nWant: %#v", have.Tag(), wantTag)
	}
	if have.Name() != wantName {
		t.Errorf("\nGot: %#v;\nWant: %#v", have.Name(), wantName)
	}
}
