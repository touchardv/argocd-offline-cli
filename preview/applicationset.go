package preview

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	appsettemplate "github.com/argoproj/argo-cd/v2/applicationset/controllers/template"
	"github.com/argoproj/argo-cd/v2/applicationset/generators"
	appsetutils "github.com/argoproj/argo-cd/v2/applicationset/utils"
	argocmd "github.com/argoproj/argo-cd/v2/cmd/argocd/commands"
	cmdutil "github.com/argoproj/argo-cd/v2/cmd/util"
	argoappv1 "github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	repoapiclient "github.com/argoproj/argo-cd/v2/reposerver/apiclient"
	"github.com/argoproj/argo-cd/v2/reposerver/metrics"
	"github.com/argoproj/argo-cd/v2/reposerver/repository"
	"github.com/argoproj/argo-cd/v2/util/argo"
	"github.com/argoproj/argo-cd/v2/util/git"
	"github.com/argoproj/pkg/errors"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	log "github.com/sirupsen/logrus"
)

const (
	applicationAPIVersion = "argoproj.io/v1alpha1"
	applicationKind       = "Application"
)

var logger *log.Logger

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	// set warn log level to avoid standard argocd info logging
	os.Setenv("ARGOCD_LOG_LEVEL", "WARN")
	cmdutil.LogLevel = "WARN"
	logger = log.StandardLogger()
	logger.SetLevel(log.WarnLevel)
}

func PreviewApplications(filename string, appName string, output string) {
	apps := generateApplications(filename)
	switch output {
	case "name":
		fmt.Println("NAME")
		for _, app := range apps {
			if !shouldMatch(appName) || appName == app.Name {
				fmt.Printf("application/%s\n", app.Name)
			}
		}
	case "json", "yaml":
		if shouldMatch(appName) {
			for _, app := range apps {
				if appName == app.Name {
					app.TypeMeta.APIVersion = applicationAPIVersion
					app.TypeMeta.Kind = applicationKind
					argocmd.PrintResource(app, output)
					break
				}
			}
		} else {
			argocmd.PrintResourceList(apps, output, false)
		}
	default:
		errors.CheckError(fmt.Errorf("unknown output format: %s", output))
	}
}

func PreviewResources(filename string, appName string, output string) {
	max, err := resource.ParseQuantity("100G")
	errors.CheckError(err)
	maxValue := max.ToDec().Value()
	initConstants := repository.RepoServerInitConstants{
		HelmManifestMaxExtractedSize:      maxValue,
		HelmRegistryMaxIndexSize:          maxValue,
		MaxCombinedDirectoryManifestsSize: max,
		StreamedManifestMaxExtractedSize:  maxValue,
		StreamedManifestMaxTarSize:        maxValue,
	}

	repoService := repository.NewService(
		metrics.NewMetricsServer(),
		NewNoopCache(),
		initConstants,
		argo.NewResourceTracking(),
		git.NoopCredsStore{},
		filepath.Join(os.TempDir(), "_argocd-offline-cli"),
	)
	if err := repoService.Init(); err != nil {
		log.Fatal("failed to initialize the repo service: ", err)
	}
	apps := generateApplications(filename)
	for _, app := range apps {
		if !shouldMatch(appName) || appName == app.Name {
			response, err := repoService.GenerateManifest(context.Background(), &repoapiclient.ManifestRequest{
				ApplicationSource: app.Spec.Source,
				AppName:           app.Name,
				Namespace:         app.Spec.Destination.Namespace,
				NoCache:           true,
				Repo: &argoappv1.Repository{
					Repo:     app.Spec.Source.RepoURL,
					Username: FindRepoUsername(app.Spec.Source.RepoURL),
					Password: FindRepoPassword(app.Spec.Source.RepoURL),
				},
				ProjectName: "applications",
			})
			if err != nil {
				log.Fatal("failed to generate manifest: ", err)
			}

			resources := map[string][]unstructured.Unstructured{}
			for _, manifest := range response.Manifests {
				resource := unstructured.Unstructured{}
				err = json.Unmarshal([]byte(manifest), &resource)
				errors.CheckError(err)

				kind := strings.ToLower(resource.GetKind())
				if _, ok := resources[kind]; !ok {
					resources[kind] = make([]unstructured.Unstructured, 0)
				}
				resources[kind] = append(resources[kind], resource)
			}
			kinds := make([]string, 0)
			for kind := range resources {
				kinds = append(kinds, kind)
			}
			sort.Strings(kinds)
			switch output {
			case "name":
				printNewline := true
				for _, kind := range kinds {
					if printNewline {
						printNewline = false
					} else {
						fmt.Println()
					}
					fmt.Println("NAME")
					for _, resource := range resources[kind] {
						fmt.Printf("%s/%s\n", kind, resource.GetName())
					}
				}
			case "json", "yaml":
				for _, kind := range kinds {
					argocmd.PrintResourceList(resources[kind], output, false)
				}

			default:
				errors.CheckError(fmt.Errorf("unknown output format: %s", output))
			}
		}
	}
}

func generateApplications(filename string) []argoappv1.Application {
	appSets, err := cmdutil.ConstructApplicationSet(filename)
	if err != nil {
		log.Fatal("failed to construct ApplicationSet: ", err)
	}
	if len(appSets) > 1 {
		log.Warnf("found %d ApplicationSets, only previewing the first entry", len(appSets))
	}
	appSet := appSets[0]
	appSetGenerators := getAppSetGenerators()
	apps, _, err := appsettemplate.GenerateApplications(log.NewEntry(logger), *appSet, appSetGenerators, &appsetutils.Render{}, nil)
	if err != nil {
		log.Fatal("failed to generate Application(s): ", err)
	}
	return apps
}

func getAppSetGenerators() map[string]generators.Generator {
	terminalGenerators := map[string]generators.Generator{
		"List": generators.NewListGenerator(),
	}
	nestedGenerators := map[string]generators.Generator{
		"List":   terminalGenerators["List"],
		"Matrix": generators.NewMatrixGenerator(terminalGenerators),
		"Merge":  generators.NewMergeGenerator(terminalGenerators),
	}
	topLevelGenerators := map[string]generators.Generator{
		"List":   terminalGenerators["List"],
		"Matrix": generators.NewMatrixGenerator(nestedGenerators),
		"Merge":  generators.NewMergeGenerator(nestedGenerators),
	}

	return topLevelGenerators
}

func shouldMatch(appName string) bool {
	return len(appName) > 0
}
