package utils

import "os"

func Readfile(path string) (string, error) {
	blob, err := os.ReadFile(path)
	return string(blob), err
}
