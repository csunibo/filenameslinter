package filenameslinter

import (
	"testing"

	"github.com/csunibo/synta"
	"github.com/liamg/memoryfs"
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
var _ = rootFS.MkdirAll("prove4", 0777)
var _ = rootFS.MkdirAll("prove4/prova4", 0777)
var _ = rootFS.WriteFile("prove4/prova4/asdiw_sdfws.txt", []byte(""), 0777)
var _ = rootFS.WriteFile("prove4/prova4/__addi11223sjdK___.txt", []byte(""), 0777)
var _ = rootFS.MkdirAll("prove5", 0777)
var _ = rootFS.MkdirAll("prove5/pro__asd", 0777)
var _ = rootFS.WriteFile("prove5/pro__asd/samu.txt", []byte(""), 0777)
var _ = rootFS.WriteFile("prove5/pro__asd/luca.md", []byte(""), 0777)
var _ = rootFS.WriteFile("prove5/pro__asd/fabio123.md", []byte(""), 0777)
var _ = rootFS.WriteFile("prove5/pro__asd/ANDREA.pdf", []byte(""), 0777)

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

func TastDirRecursiveWithCorrectDirName(t *testing.T) {
	input := `word = [a-zA-Z]+[0-9]*
ext = pdf|txt|tex|md
> word.ext`

	synta := synta.MustSynta(input)
	err := CheckDir(synta, rootFS, "prove4")
	assert.NotNil(t, err)
}

func TastDirRecursiveWithNotCorrectDirName(t *testing.T) {
	input := `word = [a-zA-Z]+[0-9]*
ext = pdf|txt|tex|md
> word.ext`

	synta := synta.MustSynta(input)
	err := CheckDir(synta, rootFS, "prove5")
	assert.NotNil(t, err)
}
