package usage

import "io"

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

func NewConfiguration(source io.Reader) *Configuration {
	return nil
}
