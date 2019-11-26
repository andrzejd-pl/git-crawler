package repositories

import (
	"gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
	"gopkg.in/src-d/go-git.v4/storage"
	"io"
)

type Repository interface {
	Download(storage.Storer, billy.Filesystem, io.Writer) error
	GetStatus() (git.Status, error)
}

type gitRepository struct {
	url     string
	key     *ssh.PublicKeys
	isBare  bool
	pointer *git.Repository
}

func (r *gitRepository) GetStatus() (git.Status, error) {
	workTree, err := r.pointer.Worktree()

	if err != nil {
		return nil, err
	}

	status, err := workTree.Status()

	if err != nil {
		return nil, err
	}

	return status, nil
}

func NewGitRepository(url string, key *ssh.PublicKeys, bare bool) Repository {
	return &gitRepository{
		url:     url,
		key:     key,
		isBare:  bare,
		pointer: nil,
	}
}

func (r *gitRepository) Download(target storage.Storer, fs billy.Filesystem, logger io.Writer) error {
	repository, err := git.Clone(target, fs,
		&git.CloneOptions{
			URL:      r.url,
			Auth:     r.key,
			Progress: logger,
		})
	r.pointer = repository

	return err
}
