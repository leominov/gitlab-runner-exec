package main

import (
	"flag"
	"os"

	"github.com/sirupsen/logrus"
)

var (
	workDir = flag.String("work-dir", "./", "Working directory.")
	remote  = flag.String("remote", "origin", "Repository remote name.")
)

func realMain() int {
	flag.Parse()
	gitlabUsername := os.Getenv("GITLAB_USER")
	gitlabPassword := os.Getenv("GITLAB_PASSWORD")
	runner, err := NewRunner(*workDir, *remote, gitlabUsername, gitlabPassword)
	if err != nil {
		logrus.Error(err)
		return 1
	}
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
