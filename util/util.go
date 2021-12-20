package util

import (
	"fmt"
	"os"
	"path"

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
