package filenameslinter

import (
	"testing"

	"github.com/liamg/memoryfs"
	"github.com/csunibo/synta"
	"github.com/stretchr/testify/assert"
)

// Don't want to use storage for tests, so we use /tmp
// but it's only linux compatible.

var rootFS = memoryfs.New()

var _ = rootFS.MkdirAll("dir1", 0777)
var _ = rootFS.WriteFile("abc", []byte(""), 0777)
var _ = rootFS.WriteFile("APPUNTI.md", []byte(""), 0777)
var _ = rootFS.WriteFile("basic.txt", []byte(""), 0777)
var _ = rootFS.MkdirAll("prove", 0777)
var _ = rootFS.WriteFile("prove/basic123.txt", []byte(""), 0777)
var _ = rootFS.MkdirAll("prove2", 0777)
var _ = rootFS.WriteFile("prove2/a1a.txt", []byte(""), 0777)
var _ = rootFS.MkdirAll("prove3", 0777)
var _ = rootFS.WriteFile("prove3/basi___c123.txt", []byte(""), 0777)
var _ = rootFS.WriteFile("prove3/pippo.txt", []byte(""), 0777)

func TestDirEmpty(t *testing.T) {
	synta := synta.MustSynta(`useless = a|dir1
> useless.useless`)
	err := CheckDir(synta, rootFS, "dir1")
	assert.Nil(t, err)
}

func TestDirOneFileCorrect(t *testing.T) {
    input := `word = [a-zA-Z]+[0-9]*
ext = pdf|txt|tex|md
> word.ext`

    synta := synta.MustSynta(input)
    err := CheckDir(synta, rootFS, "prove")
    assert.Nil(t, err)
}

func TestDirOneFileNotCorrect(t *testing.T) {
    input := `word = [a-zA-Z]+[0-9]*
ext = pdf|txt|tex|md
> word.ext`

    synta := synta.MustSynta(input)
    err := CheckDir(synta, rootFS, "prove2")
    assert.NotNil(t, err)
}

func TestDirFileCorrectAndNotCorrect(t *testing.T) {
    input := `word = [a-zA-Z]+[0-9]*
ext = pdf|txt|tex|md
> word.ext`

    synta := synta.MustSynta(input)
    err := CheckDir(synta, rootFS, "prove3")
    assert.NotNil(t, err)
    matchErr := err.(RegexMatchError)
    assert.Equal(t, "basi___c123.txt", matchErr.Filename)
}
