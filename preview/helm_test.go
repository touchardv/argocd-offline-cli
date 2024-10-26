package preview

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFindRepoCredentialsFromEnv(t *testing.T) {
	t.Setenv("HELM_REPOSITORY_CONFIG", "../testdata/repositories.yaml")
	t.Setenv("HELM_REPO_PASSWORD", "myPassword")
	t.Setenv("HELM_REPO_USERNAME", "myUsername")
	LoadLocalHelmFile()

	require.Equal(t, "myPassword", FindRepoPassword("https://dummy"))
	require.Equal(t, "myUsername", FindRepoUsername("https://dummy"))
	require.Equal(t, "myPassword", FindRepoPassword("https://unknown"))
	require.Equal(t, "myUsername", FindRepoUsername("https://unknown"))
}

func TestFindRepoCredentialsFromHelmConfig(t *testing.T) {
	t.Setenv("HELM_REPOSITORY_CONFIG", "../testdata/repositories.yaml")
	LoadLocalHelmFile()

	require.Equal(t, "helmPassword", FindRepoPassword("https://dummy"))
	require.Equal(t, "helmUsername", FindRepoUsername("https://dummy"))
	require.Equal(t, "helmPassword", FindRepoPassword("https://dummy/"))
	require.Equal(t, "helmUsername", FindRepoUsername("https://dummy/"))
	require.Equal(t, "", FindRepoPassword("https://no.defined.in.repositories"))
	require.Equal(t, "", FindRepoUsername("https://no.defined.in.repositories"))
}

func TestFindRepoCredentialsNone(t *testing.T) {
	t.Setenv("HELM_REPOSITORY_CONFIG", "no/repositories.yaml")
	LoadLocalHelmFile()

	require.Equal(t, "", FindRepoPassword("https://dummy"))
	require.Equal(t, "", FindRepoUsername("https://dummy"))
}
