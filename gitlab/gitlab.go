package gitlab

import (
	"strconv"

	"github.com/xanzy/go-gitlab"
)

type Client struct {
	cli *gitlab.Client
}

func NewClient(endpoint, username, password string) (*Client, error) {
	gitlabCli, err := gitlab.NewBasicAuthClient(nil, endpoint, username, password)
	if err != nil {
		return nil, err
	}
	return &Client{
		cli: gitlabCli,
	}, nil
}

func (c *Client) withPagination(fetch func(opts gitlab.ListOptions) (*gitlab.Response, error)) error {
	page := 1
	for {
		opts := gitlab.ListOptions{
			Page: page,
		}
		r, err := fetch(opts)
		if err != nil {
			return err
		}
		nextPageRaw := r.Header.Get("X-Next-Page")
		if len(nextPageRaw) == 0 {
			break
		}
		nextPage, err := strconv.Atoi(nextPageRaw)
		if err != nil {
			break
		}
		page = nextPage
	}
	return nil
}

func (c *Client) GetGroupVariables(gid interface{}) (map[string]string, error) {
	entries := make(map[string]string)
	err := c.withPagination(func(opts gitlab.ListOptions) (*gitlab.Response, error) {
		options := gitlab.ListGroupVariablesOptions(opts)
		fetchedEntries, r, err := c.cli.GroupVariables.ListVariables(gid, &options, nil)
		if err != nil {
			return nil, err
		}
		for _, entry := range fetchedEntries {
			entries[entry.Key] = entry.Value
		}
		return r, nil
	})
	return entries, err
}

func (c *Client) GetProjectVariables(pid interface{}) (map[string]string, error) {
	entries := make(map[string]string)
	err := c.withPagination(func(opts gitlab.ListOptions) (*gitlab.Response, error) {
		options := gitlab.ListProjectVariablesOptions(opts)
		fetchedEntries, r, err := c.cli.ProjectVariables.ListVariables(pid, &options, nil)
		if err != nil {
			return nil, err
		}
		for _, entry := range fetchedEntries {
			entries[entry.Key] = entry.Value
		}
		return r, nil
	})
	return entries, err
}
