package verutil

import (
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/tsaikd/KDGoLib/jsonex"
	"github.com/tsaikd/KDGoLib/version"
)

// GetVersionFromSource return version from source code, not from compiled binary,
func GetVersionFromSource(projectRootDir string, packageFilePath string) (ver version.Version, err error) {
	jsonBody, err := ioutil.ReadFile(packageFilePath)
	if err != nil {
		return
	}

	pkgver := struct {
		Version string `json:"version"`
	}{}

	if err = jsonex.Unmarshal(jsonBody, &pkgver); err != nil {
		return
	}

	ver.Version = pkgver.Version

	ver.GoVersion = runtime.Version()

	commitTime, err := getGitCommitTime(projectRootDir)
	if err != nil {
		return
	}
	ver.BuildTime = commitTime

	hash, err := getGitHash(projectRootDir)
	if err != nil {
		return
	}
	ver.GitCommit = hash

	return
}

// getGitCommitTime return project current git commit time
func getGitCommitTime(projroot string) (hash string, err error) {
	gitdir := filepath.Join(projroot, ".git")
	cmd := exec.Command("git", "--git-dir", gitdir, "log", "-1", "--format=%cd", "HEAD")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return
	}
	return strings.TrimSpace(string(output)), nil
}

// getGitHash return project current git hash string
func getGitHash(projroot string) (hash string, err error) {
	gitdir := filepath.Join(projroot, ".git")
	cmd := exec.Command("git", "--git-dir", gitdir, "rev-parse", "--short", "HEAD")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return
	}
	return strings.TrimSpace(string(output)), nil
}
