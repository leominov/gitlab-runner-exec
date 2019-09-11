package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/leominov/gitlab-runner-exec/git"
	"github.com/xanzy/go-gitlab"
)

const (
	gitlabCIFile = ".gitlab-ci.yml"
)

var (
	ciFilename = flag.String("ci", gitlabCIFile, "Gitlab CI configuration file.")
	workDir    = flag.String("work-dir", "./", "Working directory.")
	remote     = flag.String("remote", "origin", "Repository remote name.")
)

func main() {
	flag.Parse()
	gitCli, err := git.NewClient(*workDir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	remote, err := gitCli.Remote(*remote)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	u, err := url.Parse(remote)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	tempDir := os.TempDir()
	defer os.RemoveAll(tempDir)
	src := path.Join(*workDir, *ciFilename)
	dst := path.Join(tempDir, gitlabCIFile)
	err = CopyFile(src, dst)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	endpoint := fmt.Sprintf("%s://%s", u.Scheme, u.Hostname())
	namespace := strings.TrimPrefix(u.Path, "/")
	namespace = strings.TrimSuffix(namespace, ".git")
	gitlabCli := gitlab.NewClient(nil, os.Getenv("GITLAB_TOKEN"))
	err = gitlabCli.SetBaseURL(endpoint)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	opts := &gitlab.GetProjectOptions{}
	project, _, err := gitlabCli.Projects.GetProject(namespace, opts, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(project.Name)
	fmt.Println(project.Namespace.FullPath)
}
