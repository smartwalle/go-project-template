package go_project_template

import "fmt"

var (
	GoVersion string
	GitBranch string
	GitCommit string
	BuildTime string
)

//go build -ldflags "-X 'nx.GoVersion=`go version`' -X 'nx.GitBranch=`git rev-parse --abbrev-ref HEAD`' -X 'nx.GitCommit=`git rev-parse HEAD`' -X 'nx.BuildTime=`date "+%Z-%m-%d %H:%M:%S"`'"

func Version() string {
	return fmt.Sprintf("Go Version: %s\nGit Branch: %s\nGit Commit: %s\nBuild Time: %s\n", GoVersion, GitBranch, GitCommit, BuildTime)
}
