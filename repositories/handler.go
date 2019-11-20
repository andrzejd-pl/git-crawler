package repositories

import (
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/protocol/packp/sideband"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

const (
	isBare             bool   = false
	defaultGitUser     string = "git"
	defaultKeyPassword string = ""
	defaultKeyPath     string = "/Users/andrzejdybowski/.ssh/id_rsa"
)

type RepositoryHandler struct {
	publicKey  *ssh.PublicKeys
	url        string
	path       string
	logger     sideband.Progress
	repository *git.Repository
}

func NewRepositoryHandler(publicKey *ssh.PublicKeys, url string, path string, logger sideband.Progress) RepositoryHandler {
	return RepositoryHandler{
		publicKey: publicKey,
		url:       url,
		path:      path,
		logger:    logger,
	}
}

func NewPublicKey(keyFile, user, password string) (*ssh.PublicKeys, error) {
	if user == "" {
		user = defaultGitUser
	}

	if password == "" {
		password = defaultKeyPassword
	}

	if keyFile == "" {
		keyFile = defaultKeyPath
	}

	return ssh.NewPublicKeysFromFile(user, keyFile, password)
}

func (r RepositoryHandler) DownloadRepository() error {
	repository, err := git.PlainClone(
		r.path,
		isBare,
		&git.CloneOptions{
			URL:      r.url,
			Auth:     r.publicKey,
			Progress: r.logger,
		})

	if err != nil {
		return err
	}

	r.repository = repository
	return nil
}
