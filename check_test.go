package filenameslinter

import (
	"os"
	"testing"

	"path/filepath"

	"github.com/csunibo/synta"

	"github.com/stretchr/testify/assert"
)

// Don't want to use storage for tests, so we use /tmp
// but it's only linux compatible.
var rootDir = "/tmp/dir_test"

func TestBasicError(t *testing.T) {
	synta := synta.MustSynta("useless = a|b")
	err := CheckFilePath(synta, "abc")
	assert.NotNil(t, err)
}

func TestDirBasic(t *testing.T) {
	synta := synta.MustSynta("useless = a|b")
	err := CheckFilePath(synta, "cmd")
	assert.Nil(t, err)
}

func TestFileBasicError(t *testing.T) {
	input := `word = [a-z]+
ext = pdf|txt|tex|md
> word.ext`
	synta := synta.MustSynta(input)
	err := CheckFilePath(synta, "README.md") // does not match
	assert.NotNil(t, err)
}

func TestFileBasicNotMatch(t *testing.T) {
	input := `word = [a-z]+
ext = pdf|txt|tex|md
> word.ext`
	synta := synta.MustSynta(input)
	err := CheckFilePath(synta, "test.md") // does not exist
	assert.NotNil(t, err)
}

func TestFileBasicMatch(t *testing.T) {
	input := `word = [a-z]+
	ext = pdf|txt|tex|go
	> word.ext`

	synta := synta.MustSynta(input)
	err := CheckFilePath(synta, "check.go") // exists, and maches
	assert.Nil(t, err)
}

// DOWN THERE WE TEST DIRECTORIES:

// TODO: move this part of code to separate file, if needed in other tests.
func CreateForEach(setUp func(), tearDown func()) func(func()) {
	return func(testFunc func()) {
		setUp()
		testFunc()
		tearDown()
	}
}

var RunTest = CreateForEach(setUp, tearDown)

func setUp() {
	// create a directory structure
	// 	- dir_test
	// 		- file.txt
	// 		- file.md
	// 		- FILE.md
	// 		- FILE2.md
	// 		- file1.pdf
	// 		- dir2
	// 			- file3.txt
	// 			- file.txt

	err := os.Mkdir(rootDir, os.ModePerm)
	if err != nil {
		panic(err)
		// if the dir_test does not exist, we assume the dir structure is not created
	}

	dir2 := filepath.Join(rootDir, "dir2")
	os.Mkdir(dir2, os.ModePerm)

	dir3 := filepath.Join(rootDir, "dir3")
	os.Mkdir(dir3, os.ModePerm)

	file1 := filepath.Join(rootDir, "file.txt")
	os.Create(file1)

	file2 := filepath.Join(rootDir, "file.md")
	os.Create(file2)

	file3 := filepath.Join(rootDir, "FILE.md")
	os.Create(file3)

	file4 := filepath.Join(rootDir, "FILE2.md")
	os.Create(file4)

	file5 := filepath.Join(rootDir, "file1.pdf")
	os.Create(file5)

	file6 := filepath.Join(dir2, "file3.txt")
	os.Create(file6)

	file7 := filepath.Join(dir2, "file.txt")
	os.Create(file7)
}

func tearDown() {
	os.RemoveAll(rootDir)
}

func TestDirBasicError(t *testing.T) {
	RunTest(func() {
		input := `word = [a-z]+
		ext = pdf|txt|tex|go
		> word.ext`

		synta := synta.MustSynta(input)
		err := CheckFilePath(synta, rootDir)
		assert.NotNil(t, err)
	})
}

func TestDirBasicMatch(t *testing.T) {
	RunTest(func() {
		input := `word = [a-zA-Z]+[0-9]*
		ext = pdf|txt|tex|md
		> word.ext`

		synta := synta.MustSynta(input)
		err := CheckFilePath(synta, rootDir)
		assert.Nil(t, err)
	})
}

func TestDirDeepError(t *testing.T) {
	RunTest(func() {
		input := `word = [a-z]+
		ext = pdf|txt|tex|md
		> word.ext`

		synta := synta.MustSynta(input)
		err := CheckFilePath(synta, filepath.Join(rootDir, "dir2"))
		assert.NotNil(t, err)
	})
}

func TestDirDeepMatch(t *testing.T) {
	RunTest(func() {
		input := `word = [a-z]+[3]?
		ext = pdf|txt|tex|md
		> word.ext`

		synta := synta.MustSynta(input)
		err := CheckFilePath(synta, filepath.Join(rootDir, "dir2"))
		assert.Nil(t, err)
	})
}
