package preview

import (
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/repo"
)

var localHelmFile *repo.File

func init() {
	LoadLocalHelmFile()
}

func LoadLocalHelmFile() {
	file, err := repo.LoadFile(cli.New().RepositoryConfig)
	if err != nil {
		log.Warn("could not read helm local repository config: ", err)
	}
	localHelmFile = file
}

func FindRepoPassword(repoURL string) string {
	v, present := os.LookupEnv("HELM_REPO_PASSWORD")
	if present && strings.TrimSpace(v) != "" {
		return v
	}
	return findHelmRepo(repoURL).Password
}

func FindRepoUsername(repoURL string) string {
	v, present := os.LookupEnv("HELM_REPO_USERNAME")
	if present && strings.TrimSpace(v) != "" {
		return v
	}
	return findHelmRepo(repoURL).Username
}

func findHelmRepo(repoURL string) *repo.Entry {
	if localHelmFile != nil {
		url := strings.TrimSuffix(repoURL, "/")
		for _, r := range localHelmFile.Repositories {
			if strings.TrimSuffix(r.URL, "/") == url {
				return r
			}
		}
	}
	return &repo.Entry{}
}
