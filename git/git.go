package git

import (
	"fmt"
	"os/exec"
	"strings"
)

type Client struct {
	git string
	dir string
}

func NewClient(dir string) (*Client, error) {
	g, err := exec.LookPath("git")
	if err != nil {
		return nil, err
	}
	client := &Client{
		git: g,
		dir: dir,
	}
	return client, nil
}

func (c *Client) gitCommand(arg ...string) *exec.Cmd {
	cmd := exec.Command(c.git, arg...)
	cmd.Dir = c.dir
	return cmd
}

func (c *Client) Remote(remote string) (string, error) {
	co := c.gitCommand("remote", "get-url", remote)
	b, err := co.CombinedOutput()
	body := strings.TrimSpace(string(b))
	if err != nil {
		return "", fmt.Errorf("error getting remove url %s: %v. output: %s", remote, err, body)
	}
	return body, nil
}
