package log

import (
	"fmt"
	"log"
	"os"
)

type Channels struct {
    Info *log.Logger
	Err *log.Logger
}

func New(component string) Channels {
	flags := log.Ldate | log.Ltime
	return Channels{
		Info: log.New(os.Stdout, fmt.Sprintf("INFO %s: ", component), flags),
		Err: log.New(os.Stderr, fmt.Sprintf("ERR %s: ", component), flags),
	}
}
