package filenameslinter

import (
    "fmt"
    "os"

    "github.com/csunibo/synta"
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
    }

    return
}
