package filenameslinter

import (
	"fmt"
)

type RegexMatchError struct {
	Regexp   string
	Filename string
}

func (e RegexMatchError) Error() string {
	return fmt.Sprintf("Filname %s doesn't match the regexp %s", e.Filename, e.Regexp)
}
