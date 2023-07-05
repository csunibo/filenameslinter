package filenameslinter

import (
	"fmt"
	"regexp"
)

type RegexMatchError struct {
	Regexp   *regexp.Regexp
	Filename string
}

func (e RegexMatchError) Error() string {
	return fmt.Sprintf("Filname %s doesn't match the regexp %s", e.Filename, e.Regexp.String())
}
