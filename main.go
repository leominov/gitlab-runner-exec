package main

import (
	"flag"
	"os"

	"github.com/sirupsen/logrus"
)

var (
	ciFilename = flag.String("ci", gitlabCIFile, "Gitlab CI configuration file.")
	workDir    = flag.String("work-dir", "./", "Working directory.")
	remote     = flag.String("remote", "origin", "Repository remote name.")
)

func realMain() int {
	flag.Parse()
	runner, err := NewRunner(*ciFilename, *workDir, *remote)
	if err != nil {
		logrus.Error(err)
		return 1
	}
	defer runner.Close()
	err = runner.Exec(flag.Args())
	if err != nil {
		logrus.Error(err)
		return 1
	}
	return 0
}

func main() {
	os.Exit(realMain())
}
