package main

import (
	"fmt"
	"github.com/andrzejd-pl/git-crawler/csv"
	"github.com/andrzejd-pl/git-crawler/html"
	"github.com/andrzejd-pl/git-crawler/repositories"
	"github.com/andrzejd-pl/git-crawler/usage"
	"os"
	"strconv"
	"sync"
)

const (
	standardPath string = "app/views/parts/site-footer.blade.php"
)

func main() {
	logger := os.Stderr
	pathToMainDir := os.Args[1]
	maxThreads, err := strconv.Atoi(os.Args[2])
	usage.CheckErrorWithPanic(os.Stderr, err)

	sitesFile, err := os.Open(os.Args[3])
	usage.CheckErrorWithPanic(logger, err)

	sites, err := csv.ReadSites(sitesFile)
	usage.CheckErrorWithPanic(logger, err)

	publicKey, err := repositories.NewPublicKey("", "", "")
	usage.CheckError(os.Stdout, err, true)

	var wg sync.WaitGroup
	i := 0

	for siteId, urlRepo := range sites {
		wg.Add(1)
		i++

		go func(id, url string) {
			defer wg.Done()
			pathToRepo := pathToMainDir + "/" + id

			fileToLog, err := os.Create("./logs/" + id + ".log")
			usage.CheckErrorWithOnlyLogging(logger, err)

			repositoryHandler := repositories.NewRepositoryHandler(publicKey, url, pathToRepo, fileToLog)
			err = repositoryHandler.DownloadRepository()
			usage.CheckErrorWithOnlyLogging(logger, err)

			err = replace(pathToRepo + "/" + standardPath)
			usage.CheckErrorWithOnlyLogging(logger, err)

		}(siteId, urlRepo)

		if i%maxThreads == 0 {
			wg.Wait()
		}
	}

	wg.Wait()
	fmt.Println("ok")
}

func replace(filePath string) error {
	tempFileName := filePath + ".temp"

	sourceFile, err := os.Open(filePath)
	defer sourceFile.Close()

	if err != nil {
		return err
	}

	targetFile, err := os.Create(tempFileName)
	defer targetFile.Close()

	if err != nil {
		return err
	}

	err = html.Replace(sourceFile, targetFile)

	if err != nil {
		return err
	}

	err = sourceFile.Close()

	if err != nil {
		return err
	}
	err = targetFile.Close()

	if err != nil {
		return err
	}

	err = os.Rename(tempFileName, filePath)

	if err != nil {
		return err
	}

	return nil
}
