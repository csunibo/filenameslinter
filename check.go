package filenameslinter

import (
    "fmt"
    "os"

    "github.com/csunibo/synta"
    "github.com/csunibo/synta/regexp"
)

func CheckFileName(synta synta.Synta, path string) (err error) {
    file, err := os.Open(path)
    if err != nil {
        err = fmt.Errorf("File %s does not exists: %v", path, err)
        return
    }

    info, err := file.Stat()
    if err != nil {
        err = fmt.Errorf("Error while getting stats for file %s: %v", path, err)
        return
    }

    if info.IsDir() {
        fmt.Println("It's a directory")
        // TODO handling for directory
        return
    }

    reg, err := regexp.Convert(synta)
    if err != nil {
        err = fmt.Errorf("Can't convert synta to regexp, %v", err)
    }

    if reg.Match([]byte(info.Name())) {
        fmt.Println("Match ok")
    } else {
        err = fmt.Errorf("Regexp don't match, regexp: %s, filename: %s", reg.String(), info.Name())
    }

    return
}
