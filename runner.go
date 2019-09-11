package main

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/leominov/gitlab-runner-exec/git"
	"github.com/leominov/gitlab-runner-exec/gitlab"
	"github.com/sirupsen/logrus"
)

const (
	gitlabCIFile = ".gitlab-ci.yml"
)

type Runner struct {
	gitCli    *git.Client
	gitlabCli *gitlab.Client
	ciFile    string
	workDir   string
	tempDir   string
	remote    string
	namespace string
	groups    []string
}

func NewRunner(ciFile, wd, remote string) (*Runner, error) {
	gitCli, err := git.NewClient(*workDir)
	if err != nil {
		return nil, err
	}
	r := &Runner{
		gitCli:  gitCli,
		ciFile:  ciFile,
		workDir: wd,
		remote:  remote,
		tempDir: os.TempDir(),
	}
	endpoint, namespace, err := r.parseRemote()
	if err != nil {
		return nil, err
	}
	r.namespace = namespace
	r.groups = GroupsFromNamespace(namespace)
	gitlabCli, err := gitlab.NewClient(endpoint, os.Getenv("GITLAB_TOKEN"))
	if err != nil {
		return nil, err
	}
	r.gitlabCli = gitlabCli
	err = r.prepareWorkspace()
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (r *Runner) Close() error {
	if len(r.tempDir) == 0 {
		return nil
	}
	return os.RemoveAll(r.tempDir)
}

func (r *Runner) prepareWorkspace() error {
	src := path.Join(r.workDir, r.ciFile)
	dst := path.Join(r.tempDir, gitlabCIFile)
	err := CopyFile(src, dst)
	if err != nil {
		return err
	}
	return nil
}

func (r *Runner) parseRemote() (endpoint string, namespace string, err error) {
	gitCli, err := git.NewClient(*workDir)
	if err != nil {
		return
	}
	rm, err := gitCli.Remote(r.remote)
	if err != nil {
		return
	}
	u, err := url.Parse(rm)
	if err != nil {
		return
	}
	endpoint = fmt.Sprintf("%s://%s", u.Scheme, u.Hostname())
	namespace = strings.TrimPrefix(u.Path, "/")
	namespace = strings.TrimSuffix(namespace, ".git")
	return
}

func (r *Runner) getVariables() (map[string]string, error) {
	vars := make(map[string]string)
	for _, group := range r.groups {
		vs, err := r.gitlabCli.GetGroupVariables(group)
		if err != nil {
			return vars, err
		}
		for k, v := range vs {
			vars[k] = v
		}
	}
	projectVars, err := r.gitlabCli.GetProjectVariables(r.namespace)
	if err != nil {
		return vars, err
	}
	for k, v := range projectVars {
		vars[k] = v
	}
	return vars, nil
}

func (r *Runner) Exec(userArgs []string) error {
	vars, err := r.getVariables()
	if err != nil {
		return err
	}
	args := []string{"exec"}
	args = append(args, userArgs...)
	for k, v := range vars {
		envVar := fmt.Sprintf("--env=%s=%q", k, v)
		args = append(args, envVar)
	}
	logrus.WithField("src", "gitlab").Infof("Found variables: %d", len(vars))
	cmd := exec.Command("gitlab-runner", args...)
	cmd.Stderr = logrus.WithField("src", "cmd").WriterLevel(logrus.ErrorLevel)
	cmd.Stdout = logrus.WithField("src", "cmd").WriterLevel(logrus.InfoLevel)
	return cmd.Run()
}
