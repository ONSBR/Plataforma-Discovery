package env

import (
	"fmt"
)

//GetGitServerKeysPath returns git server keys path
func GetGitServerKeysPath() string {
	return Get("GIT_SERVER_KEYS_PATH", fmt.Sprintf("/git-server/keys"))
}

//GetGitServerReposPath returns git server repos path
func GetGitServerReposPath() string {
	return Get("GIT_SERVER_REPOS_PATH", fmt.Sprintf("/git-server/repos"))
}

//GetDeploysPath returns where artifacts will be cloned to be deployed
func GetDeploysPath() string {
	return Get("DEPLOY_PATH", fmt.Sprintf("worker/deploys"))
}

//GetSSHRemoteURL returns git remote url pattern for ssh protocol
func GetSSHRemoteURL(solution, app string) string {
	user := Get("GIT_SERVER_USER", "git")
	host := Get("GET_SERVER_EXTERNAL_HOST", "localhost")
	port := Get("GET_SERVER_PORT", "2222")
	return fmt.Sprintf("ssh://%s@%s:%s/git-server/repos/%s/%s", user, host, port, solution, app)
}

//GetSSHRemoteInternalURL returns git remote url pattern for ssh protocol for internal access
func GetSSHRemoteInternalURL(solution, app string) string {
	user := Get("GIT_SERVER_USER", "git")
	host := Get("GET_SERVER_HOST", "git-server")
	port := Get("GET_SERVER_PORT", "2222")
	return fmt.Sprintf("ssh://%s@%s:%s/git-server/repos/%s/%s", user, host, port, solution, app)
}
