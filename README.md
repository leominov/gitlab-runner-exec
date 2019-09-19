# gitlab-runner-exec

Runs task locally with all of project variables from GitLab.

## Usage

Before run check access for project and groups variables.

```shell
export GITLAB_USER=USER
export GITLAB_PASSWORD=PASSWORD
gitlab-runner-exec --env=FOO=BAR docker lint
```
