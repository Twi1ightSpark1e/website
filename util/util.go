package util

import (
	"fmt"
	"log"
	"os"
	"path"
)

func ExecPath() string {
	exec, err := os.Executable()
	if (err != nil) {
		log.Fatal(err)
	}

	return path.Dir(exec)
}

func BasePath(suffix string) string {
	return fmt.Sprintf("%s/%s", ExecPath(), suffix)
}

