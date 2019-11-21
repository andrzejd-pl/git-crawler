package html

import (
	"bufio"
	"regexp"
)

type Replacer interface {
	Replace(*bufio.Scanner, *bufio.Writer) error
}

type replacer struct {
	searchPattern  string
	replacePattern string
}

func NewReplacer(search, replace string) Replacer {
	return replacer{
		searchPattern:  search,
		replacePattern: replace,
	}
}

func (r replacer) Replace(source *bufio.Scanner, target *bufio.Writer) error {
	regexEngine := regexp.MustCompile(r.searchPattern)

	for source.Scan() {
		text := source.Text()
		var err error

		if regexEngine.MatchString(text) {
			_, err = target.WriteString(regexEngine.ReplaceAllString(text, r.replacePattern) + "\n")
		} else {
			_, err = target.WriteString(text + "\n")
		}

		if err != nil {
			return err
		}

		if err := target.Flush(); err != nil {
			return err
		}
	}

	return source.Err()
}
