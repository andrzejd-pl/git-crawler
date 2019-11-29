package usage

import (
	"strings"
	"testing"
)

func TestNewConfiguration(t *testing.T) {
	newConfigurationTests := []struct {
		input string
		want  Configuration
	}{
		{
			"",
			Configuration{
				StandardFilePath: "",
				DefaultKeyPath:   "",
				SearchPattern:    "",
				ReplacePattern:   "",
				TempExtension:    "",
				NewBranch:        "",
				CommitMessage:    "",
				AuthorName:       "",
				AuthorEmail:      "",
			},
		},
		{
			"keyPath: ./.ssh/id_rsa\n" +
				"fileToReplace: html/footer.go\n" +
				"SearchPattern: (www)?\\.?wp\\.pl\n" +
				"ReplacePattern: onet.pl\n" +
				"temporaryExtension: .xd\n" +
				"newBranchName: feature/tests\n" +
				"CommitMessage: \"test(test): test\"\n" +
				"AuthorName: Jan Kowalski\n" +
				"AuthorEmail: jan@kowalski.pl\n",
			Configuration{
				StandardFilePath: "html/footer.go",
				DefaultKeyPath:   "./.ssh/id_rsa",
				SearchPattern:    "(www)?\\.?wp\\.pl",
				ReplacePattern:   "onet.pl",
				TempExtension:    ".xd",
				NewBranch:        "feature/tests",
				CommitMessage:    "test(test): test",
				AuthorName:       "Jan Kowalski",
				AuthorEmail:      "jan@kowalski.pl",
			},
		},
	}

	for _, tt := range newConfigurationTests {
		configuration, err := NewConfiguration(strings.NewReader(tt.input))

		if err != nil {
			t.Error(err.Error())
			continue
		}

		if *configuration != tt.want {
			t.Errorf("got %v want %v", *configuration, tt.want)
		}
	}
}
