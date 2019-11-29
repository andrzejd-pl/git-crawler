package usage

import (
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
)

type Configuration struct {
	DefaultKeyPath   string
	StandardFilePath string
	SearchPattern    string
	ReplacePattern   string
	TempExtension    string
	NewBranch        string
	CommitMessage    string
	AuthorName       string
	AuthorEmail      string
}

func NewConfiguration(source io.Reader) (*Configuration, error) {
	var yamlConfig struct {
		DefaultKeyPath   string `yaml:"keyPath"`
		StandardFilePath string `yaml:"fileToReplace"`
		SearchPattern    string `yaml:"SearchPattern"`
		ReplacePattern   string `yaml:"ReplacePattern"`
		TempExtension    string `yaml:"temporaryExtension"`
		NewBranch        string `yaml:"newBranchName"`
		CommitMessage    string `yaml:"CommitMessage"`
		AuthorName       string `yaml:"AuthorName"`
		AuthorEmail      string `yaml:"AuthorEmail"`
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
		DefaultKeyPath:   yamlConfig.DefaultKeyPath,
		StandardFilePath: yamlConfig.StandardFilePath,
		SearchPattern:    yamlConfig.SearchPattern,
		ReplacePattern:   yamlConfig.ReplacePattern,
		TempExtension:    yamlConfig.TempExtension,
		NewBranch:        yamlConfig.NewBranch,
		CommitMessage:    yamlConfig.CommitMessage,
		AuthorName:       yamlConfig.AuthorName,
		AuthorEmail:      yamlConfig.AuthorEmail,
	}, nil
}
