package repositories

import (
	"gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
	"gopkg.in/src-d/go-git.v4/storage"
	"io"
	"time"
)

type Repository interface {
	Download(storage.Storer, billy.Filesystem, io.Writer) error
	GetStatus() (git.Status, error)
	CommitAllChanges(string, string, string) error
	CheckoutBranch(string) error
	PushChanges(io.Writer) error
}

type gitRepository struct {
	url     string
	key     *ssh.PublicKeys
	isBare  bool
	pointer *git.Repository
}

func (r *gitRepository) CheckoutBranch(branchName string) error {
	workTree, err := r.pointer.Worktree()

	if err != nil {
		return err
	}

	return workTree.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName(branchName),
		Create: true,
	})
}

func (r *gitRepository) PushChanges(logger io.Writer) error {
	return r.pointer.Push(&git.PushOptions{
		Progress: logger,
	})
}

func (r *gitRepository) CommitAllChanges(commitMessage, authorName, authorEmail string) error {
	workTree, err := r.pointer.Worktree()
	var hashes []plumbing.Hash

	if err != nil {
		return err
	}

	status, err := r.GetStatus()

	if err != nil {
		return err
	}

	for fileName, _ := range status {
		hash, err := workTree.Add(fileName)

		if err != nil {
			return err
		}

		hashes = append(hashes, hash)
	}

	_, err = workTree.Commit(commitMessage, &git.CommitOptions{
		All: true,
		Author: &object.Signature{
			Name:  authorName,
			Email: authorEmail,
			When:  time.Now(),
		},
	})

	return err
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
