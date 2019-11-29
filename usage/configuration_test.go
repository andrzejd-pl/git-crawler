package usage

import (
	"strings"
	"testing"
)

func TestNewConfiguration(t *testing.T) {
	got := ""
	want := Configuration{
		standardFilePath: "",
		defaultKeyPath:   "",
		searchPattern:    "",
		replacePattern:   "",
		tempExtension:    "",
		newBranch:        "",
		commitMessage:    "",
		authorName:       "",
		authorEmail:      "",
	}

	configuration := NewConfiguraion(strings.NewReader(got))

	if configuration != want {
		t.Errorf("got %g want %g", configuration, want)
	}
}
