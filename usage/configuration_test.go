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
				MaxThreads:       0,
				RepositoriesFile: "",
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
				"authorEmail: jan@kowalski.pl\n" +
				"maxThreads: 4\n" +
				"repositoriesFile: ./storage/repos.csv\n",
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
				MaxThreads:       4,
				RepositoriesFile: "./storage/repos.csv",
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
