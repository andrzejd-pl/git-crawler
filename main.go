package main

import (
	"bufio"
	"github.com/andrzejd-pl/git-crawler/csv"
	"github.com/andrzejd-pl/git-crawler/html"
	"github.com/andrzejd-pl/git-crawler/repositories"
	"github.com/andrzejd-pl/git-crawler/usage"
	"gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-billy.v4/memfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
	"gopkg.in/src-d/go-git.v4/storage/memory"
	"os"
	"strconv"
	"sync"
)

func main() {
	logger := os.Stderr
	maxThreads, err := strconv.Atoi(os.Args[2])
	usage.CheckErrorWithPanic(os.Stderr, err)

	sitesFile, err := os.Open(os.Args[3])
	usage.CheckErrorWithPanic(logger, err)

	sites, err := csv.ReadSites(sitesFile)
	usage.CheckErrorWithPanic(logger, err)

	publicKey, err := ssh.NewPublicKeysFromFile("git", defaultKeyPath, "")
	usage.CheckError(os.Stdout, err, true)

	var wg sync.WaitGroup
	threadsNumber := 0

	for siteId, urlRepo := range sites {
		wg.Add(1)
		threadsNumber++

		go func(id, url string) {
			defer wg.Done()

			fileToLog, err := os.Create("./logs/" + id + ".log")
			usage.CheckErrorWithOnlyLogging(fileToLog, err)
			err = thread(repositories.NewGitRepository(url, publicKey, false), fileToLog)
			usage.CheckErrorWithOnlyLogging(fileToLog, err)

			if err == nil {
				_, _ = fileToLog.WriteString(id + ": ok")
			}
		}(siteId, urlRepo)

		if threadsNumber%maxThreads == 0 {
			_, _ = logger.WriteString("wait")
			wg.Wait()
		}
	}

	wg.Wait()
	_, _ = logger.WriteString("ok")
}

func thread(repo repositories.Repository, logFile *os.File) error {
	storage := memory.NewStorage()
	fileSystem := memfs.New()
	err := repo.Download(storage, fileSystem, logFile)

	if err != nil {
		return err
	}

	err = repo.CheckoutBranch(newBranch)

	if err != nil {
		return err
	}

	fileName := standardPath
	tempFileName := fileName + tempExtension
	source, target, err := openFileWithCreatingTempFile(fileName, fileSystem)

	if err != nil {
		return err
	}

	defer source.Close()
	defer target.Close()

	err = replaceAndCloseFiles(source, target)

	if err != nil {
		return err
	}

	err = fileSystem.Rename(tempFileName, fileName)

	if err != nil {
		return err
	}

	err = repo.CommitAllChanges(commitMessage, authorName, authorEmail)

	if err != nil {
		return err
	}
	err = repo.PushChanges(os.Stdout, plumbing.ReferenceName(newBranch))

	if err == git.NoErrAlreadyUpToDate {

	}

	return err
}

func openFileWithCreatingTempFile(name string, fileSystem billy.Filesystem) (billy.File, billy.File, error) {
	source, err := fileSystem.Open(name)

	if err != nil {
		return nil, nil, err
	}

	target, err := fileSystem.Create(name + tempExtension)

	if err != nil {
		return nil, nil, err
	}

	return source, target, nil
}

func replaceAndCloseFiles(sourceFile billy.File, targetFile billy.File) error {
	replacer := html.NewReplacer(patternOldLink, newLinkValue)
	err := replacer.Replace(bufio.NewScanner(sourceFile), bufio.NewWriter(targetFile))

	if err != nil {
		return err
	}

	err = sourceFile.Close()

	if err != nil {
		return err
	}

	return targetFile.Close()
}
