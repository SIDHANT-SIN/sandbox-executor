package main

import (
    "os"
)

func writeCode(dir, filename, code string) error {
    path := dir + "/" + filename
    return os.WriteFile(path, []byte(code), 0644)
}