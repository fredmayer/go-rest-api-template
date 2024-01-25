package logging

import (
	"fmt"
	"os"
	"path"
	"runtime"

	"github.com/sirupsen/logrus"
)

var (
	logger *logrus.Logger
)

func Get() *logrus.Logger {
	return logger
}

func Log() *logrus.Logger {
	return logger
}

func Init(level string) *logrus.Logger {
	l := logrus.New()

	l.SetReportCaller(true)

	l.SetFormatter(&logrus.TextFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
			filename := path.Base(f.File)
			//fmt.Sprintf("%s()", f.Func.Name())
			return "", fmt.Sprintf("%s:%d", filename, f.Line)
		},
		DisableColors: false,
		FullTimestamp: true,
	})

	err := os.MkdirAll("logs", 0755)
	if err != nil && !os.IsExist(err) {
		panic(err)
	}

	// logFile, err := os.OpenFile("logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
	// if err != nil {
	// 	panic(err)
	// }

	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		panic(err)
	}

	l.SetLevel(lvl)

	//l.SetOutput(logFile)
	logger = l

	return logger
}
