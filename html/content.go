package html

import (
	"bufio"
	"io"
	"regexp"
)

const (
	patternOldLink string = "(created-hld__link.*href\\=\\\").*(\\\")"
	newLinkValue   string = "$1{{ __('cms.created_by') }}$2"
)

func Replace(sourceFile io.Reader, destinedFile io.Writer) error {
	scanner := bufio.NewScanner(sourceFile)
	writer := bufio.NewWriter(destinedFile)
	regexEngine := regexp.MustCompile(patternOldLink)

	for scanner.Scan() {
		text := scanner.Text()
		var err error

		if regexEngine.MatchString(text) {
			_, err = writer.WriteString(regexEngine.ReplaceAllString(text, newLinkValue) + "\n")
		} else {
			_, err = writer.WriteString(text + "\n")
		}

		if err != nil {
			return err
		}

		if err := writer.Flush(); err != nil {
			return err
		}
	}

	return scanner.Err()
}
