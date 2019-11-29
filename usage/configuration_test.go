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
				standardFilePath: "",
				defaultKeyPath:   "",
				searchPattern:    "",
				replacePattern:   "",
				tempExtension:    "",
				newBranch:        "",
				commitMessage:    "",
				authorName:       "",
				authorEmail:      "",
			},
		},
		{
			"keyPath: ./.ssh/id_rsa\n" +
				"fileToReplace: html/footer.go\n" +
				"searchPattern: (www)?\\.?wp\\.pl\n" +
				"replacePattern: onet.pl\n" +
				"temporaryExtension: .xd\n" +
				"newBranchName: feature/tests\n" +
				"commitMessage: \"test(test): test\"\n" +
				"authorName: Jan Kowalski\n" +
				"authorEmail: jan@kowalski.pl\n",
			Configuration{
				standardFilePath: "html/footer.go",
				defaultKeyPath:   "./.ssh/id_rsa",
				searchPattern:    "(www)?\\.?wp\\.pl",
				replacePattern:   "onet.pl",
				tempExtension:    ".xd",
				newBranch:        "feature/tests",
				commitMessage:    "test(test): test",
				authorName:       "Jan Kowalski",
				authorEmail:      "jan@kowalski.pl",
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
