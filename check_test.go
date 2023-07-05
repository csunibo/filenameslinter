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
