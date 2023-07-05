package filenameslinter

import (
    "testing"
    "github.com/csunibo/synta"

    "github.com/stretchr/testify/assert"
)

func TestBasicError(t *testing.T) {
    synta, _ := synta.ParseSynta("")
    err := CheckFileName(synta, "abc")
    assert.NotNil(t, err)
}

func TestBasicDir(t *testing.T) {
    synta, _ := synta.ParseSynta("")
    err := CheckFileName(synta, "cmd")
    assert.Nil(t, err)
}

func TestBasicFileError(t *testing.T) {
    input := `word = [a-z]+
ext = pdf|txt|tex|md
> word.ext`
    synta, _ := synta.ParseSynta(input)
    err := CheckFileName(synta, "README.md")
    assert.NotNil(t, err)
}
