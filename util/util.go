package util

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/Twi1ightSpark1e/website/log"
)

func ExecPath() string {
	exec, err := os.Executable()
	if (err != nil) {
		logger := log.New("Util")
		logger.Err.Fatal(err)
	}

	return path.Dir(exec)
}

func FullPath(suffix string) string {
	if suffix[0] == os.PathSeparator {
		return suffix
	}
	return fmt.Sprintf("%s/%s", ExecPath(), suffix)
}

// https://stackoverflow.com/a/26809999
func Glob(dir string, ext string) ([]string, error) {
  files := []string{}
  err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
    if filepath.Ext(path) == ext {
      files = append(files, path)
    }
    return nil
  })

  return files, err
}
