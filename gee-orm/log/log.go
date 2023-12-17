package log

import (
	"log"
	"os"
)

var (
	// Adding color and style in the terminal
	// log.LstdFlags: The log message's prefix will contain the current time.
	// log.Lshortfile: The log message's prefix will contain the name of the file and the line number where the log entry was generated.
	errorLog = log.New(os.Stdout, "\033[31m[error]\033[0m ", log.LstdFlags|log.Lshortfile)
	infoLog  = log.New(os.Stdout, "\033[33m[info]\033[0m ", log.LstdFlags|log.Lshortfile)
)

var (
	Error  = errorLog.Println
	Errorf = errorLog.Printf
	Info   = infoLog.Println
	Infof  = infoLog.Printf
)

const (
	INFOLEVEL = iota
	ERRORLEVEL
	DISCARD
)

func SetLevel(level int) {
	switch level {
	case INFOLEVEL:
		{

		}
	}
}
