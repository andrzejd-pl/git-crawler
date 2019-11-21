package repositories

import (
	"gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
	"gopkg.in/src-d/go-git.v4/storage"
	"io"
)

type Repository interface {
	Download(target storage.Storer, fs billy.Filesystem, logger io.Writer) (*git.Repository, error)
}

type gitRepository struct {
	url    string
	key    *ssh.PublicKeys
	isBare bool
}

func NewGitRepository(url string, key *ssh.PublicKeys, bare bool) Repository {
	return gitRepository{
		url:    url,
		key:    key,
		isBare: bare,
	}
}

func (r gitRepository) Download(target storage.Storer, fs billy.Filesystem, logger io.Writer) (*git.Repository, error) {
	repository, err := git.Clone(target, fs,
		&git.CloneOptions{
			URL:      r.url,
			Auth:     r.key,
			Progress: logger,
		})

	if err != nil {
		return nil, err
	}

	return repository, nil
}
