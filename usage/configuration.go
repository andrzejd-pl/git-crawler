package usage

import (
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
)

type Configuration struct {
	defaultKeyPath   string
	standardFilePath string
	searchPattern    string
	replacePattern   string
	tempExtension    string
	newBranch        string
	commitMessage    string
	authorName       string
	authorEmail      string
}

func NewConfiguration(source io.Reader) (*Configuration, error) {
	var yamlConfig struct {
		DefaultKeyPath   string `yaml:"keyPath"`
		StandardFilePath string `yaml:"fileToReplace"`
		SearchPattern    string `yaml:"searchPattern"`
		ReplacePattern   string `yaml:"replacePattern"`
		TempExtension    string `yaml:"temporaryExtension"`
		NewBranch        string `yaml:"newBranchName"`
		CommitMessage    string `yaml:"commitMessage"`
		AuthorName       string `yaml:"authorName"`
		AuthorEmail      string `yaml:"authorEmail"`
	}

	yamlData, err := ioutil.ReadAll(source)

	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(yamlData, &yamlConfig)

	if err != nil {
		return nil, err
	}

	return &Configuration{
		defaultKeyPath:   yamlConfig.DefaultKeyPath,
		standardFilePath: yamlConfig.StandardFilePath,
		searchPattern:    yamlConfig.SearchPattern,
		replacePattern:   yamlConfig.ReplacePattern,
		tempExtension:    yamlConfig.TempExtension,
		newBranch:        yamlConfig.NewBranch,
		commitMessage:    yamlConfig.CommitMessage,
		authorName:       yamlConfig.AuthorName,
		authorEmail:      yamlConfig.AuthorEmail,
	}, nil
}
