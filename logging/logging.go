package logging

import (
	"log"
	"os"
)

var (
	InfoLog *log.Logger
	WarnLog *log.Logger
	ErrLog  *log.Logger
)

func init() {

	InfoLog = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarnLog = log.New(os.Stdout, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrLog = log.New(os.Stderr, "Err: ", log.Ldate|log.Ltime|log.Lshortfile)
}
