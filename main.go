package main

import (
	"bufio"
	"errors"
	"flag"
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
	"sync"
)

var configuration *usage.Configuration

func main() {
	logger := os.Stderr
	var fileName string
	flag.StringVar(&fileName, "C", "", "Configuration YAML file.")
	flag.Parse()

	if fileName == "" {
		usage.CheckErrorWithPanic(logger, errors.New("please provide configuration yaml file by using -C option"))
	}

	configurationFile, err := os.Open(fileName)
	usage.CheckErrorWithPanic(logger, err)

	configuration, err = usage.NewConfiguration(configurationFile)
	usage.CheckErrorWithPanic(logger, err)

	sitesFile, err := os.Open(configuration.RepositoriesFile)
	usage.CheckErrorWithPanic(logger, err)

	sites, err := csv.ReadSites(sitesFile)
	usage.CheckErrorWithPanic(logger, err)

	publicKey, err := ssh.NewPublicKeysFromFile("git", configuration.DefaultKeyPath, "")
	usage.CheckError(logger, err, true)

	var wg sync.WaitGroup
	threadsNumber := 0

	for siteId, urlRepo := range sites {
		wg.Add(1)
		threadsNumber++

		go func(id, url string) {
			defer wg.Done()
			fileError, err := os.Create("./logs/" + id + "-error.log")

			fileToLog, err := os.Create("./logs/" + id + "-git.log")
			usage.CheckErrorWithOnlyLogging(fileError, err)
			err = thread(repositories.NewGitRepository(url, publicKey, false), fileToLog)
			usage.CheckErrorWithOnlyLogging(fileError, err)

			if err == nil {
				_, _ = fileError.WriteString(id + ": ok\n")
			}
		}(siteId, urlRepo)

		if threadsNumber%configuration.MaxThreads == 0 {
			_, _ = logger.WriteString("wait\n")
			wg.Wait()
		}
	}

	wg.Wait()
	_, _ = logger.WriteString("ok\n")
}

func thread(repo repositories.Repository, logFile *os.File) error {
	storage := memory.NewStorage()
	fileSystem := memfs.New()
	defer fileSystem.Remove(".")
	err := repo.Download(storage, fileSystem, logFile)

	if err != nil {
		return err
	}

	err = repo.CheckoutBranch(configuration.NewBranch)

	if err != nil {
		return err
	}

	fileName := configuration.StandardFilePath
	tempFileName := fileName + configuration.TempExtension
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

	err = repo.CommitAllChanges(configuration.CommitMessage, configuration.AuthorName, configuration.AuthorEmail)

	if err != nil {
		return err
	}
	err = repo.PushChanges(os.Stdout, plumbing.ReferenceName(configuration.NewBranch))

	if err == git.NoErrAlreadyUpToDate {

	}

	return err
}

func openFileWithCreatingTempFile(name string, fileSystem billy.Filesystem) (billy.File, billy.File, error) {
	source, err := fileSystem.Open(name)

	if err != nil {
		return nil, nil, err
	}

	target, err := fileSystem.Create(name + configuration.TempExtension)

	if err != nil {
		return nil, nil, err
	}

	return source, target, nil
}

func replaceAndCloseFiles(sourceFile billy.File, targetFile billy.File) error {
	replacer := html.NewReplacer(configuration.SearchPattern, configuration.ReplacePattern)
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
