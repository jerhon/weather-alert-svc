package logging

import "log"

var format = log.Ltime | log.Lmicroseconds | log.Lshortfile

var (
	Info,
	Debug,
	Error,
	Warning,
	Fatal *log.Logger
)

func Init() {
	Info = log.New(log.Writer(), "[INFO] ", format)
	Debug = log.New(log.Writer(), "[DBG ] ", format)
	Error = log.New(log.Writer(), "[ERR ] ", format)
	Warning = log.New(log.Writer(), "[WARN] ", format)
	Fatal = log.New(log.Writer(), "[FATAL] ", format)
}
