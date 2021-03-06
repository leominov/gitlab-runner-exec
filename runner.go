package main

import (
	"fmt"
	"net/url"
	"os/exec"
	"strings"

	"github.com/leominov/gitlab-runner-exec/git"
	"github.com/leominov/gitlab-runner-exec/gitlab"
	"github.com/sirupsen/logrus"
)

type Runner struct {
	gitCli       *git.Client
	gitlabCli    *gitlab.Client
	runnerBinary string
	workDir      string
	remote       string
	namespace    string
	groups       []string
}

func NewRunner(wd, remote, username, password string) (*Runner, error) {
	runnerBinary, err := exec.LookPath("gitlab-runner")
	if err != nil {
		return nil, err
	}
	gitCli, err := git.NewClient(*workDir)
	if err != nil {
		return nil, err
	}
	r := &Runner{
		runnerBinary: runnerBinary,
		gitCli:       gitCli,
		workDir:      wd,
		remote:       remote,
	}
	endpoint, namespace, err := r.parseRemote()
	if err != nil {
		return nil, err
	}
	r.namespace = namespace
	r.groups = GroupsFromNamespace(namespace)
	gitlabCli, err := gitlab.NewClient(endpoint, username, password)
	if err != nil {
		return nil, err
	}
	r.gitlabCli = gitlabCli
	return r, nil
}

func (r *Runner) parseRemote() (endpoint string, namespace string, err error) {
	rm, err := r.gitCli.Remote(r.remote)
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

func (r *Runner) Exec(userArgs []string, envs map[string]string) error {
	vars, err := r.getVariables()
	if err != nil {
		return err
	}
	args := []string{"exec"}
	args = append(args, userArgs...)
	for k, v := range envs {
		vars[k] = v
	}
	for k, v := range vars {
		env := fmt.Sprintf("%s=%s", k, v)
		args = append(args, "--env", env)
	}
	logrus.WithField("src", "gitlab").Infof("Found variables: %d", len(vars))
	cmd := exec.Command(r.runnerBinary, args...)
	cmd.Dir = r.workDir
	cmd.Stderr = logrus.WithField("src", "cmd").WriterLevel(logrus.ErrorLevel)
	cmd.Stdout = logrus.WithField("src", "cmd").WriterLevel(logrus.InfoLevel)
	return cmd.Run()
}
