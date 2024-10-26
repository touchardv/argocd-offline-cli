package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/touchardv/argocd-offline-cli/preview"
)

func NewCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "argocd-offline-cli",
		Short: "An Argo CD CLI offline utility",
		Long: `A utility, based on Argo CD, that can be used "offline" (without requiring a running Argo CD server),
to preview the Kubernetes resource manifests being created and managed by Argo CD.`,
	}

	rootCmd.AddCommand(AppSetCommand())

	return rootCmd
}

func AppSetCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "appset",
		Short: "Preview ApplicationSets",
	}
	command.AddCommand(PreviewApplicationsCommand())
	command.AddCommand(PreviewApplicationResourcesCommand())
	return command
}

func PreviewApplicationsCommand() *cobra.Command {
	var name string
	var output string
	command := &cobra.Command{
		Use:   "preview-apps APPSETMANIFEST",
		Short: "Preview Application(s) generated from an ApplicationSet",
		Run: func(c *cobra.Command, args []string) {
			if len(args) == 0 {
				c.HelpFunc()
				os.Exit(1)
			}
			filename := args[0]
			preview.PreviewApplications(filename, name, output)
		},
	}
	command.Flags().StringVarP(&name, "name", "n", "", "Name of the Application to preview")
	command.Flags().StringVarP(&output, "output", "o", "name", "Output format. One of: name|json|yaml")
	return command
}

func PreviewApplicationResourcesCommand() *cobra.Command {
	var name string
	var output string
	command := &cobra.Command{
		Use:   "preview-resources APPSETMANIFEST",
		Short: "Preview Kubernetes resource(s) generated from an ApplicationSet/Application",
		Run: func(c *cobra.Command, args []string) {
			if len(args) == 0 {
				c.HelpFunc()
				os.Exit(1)
			}
			filename := args[0]
			preview.PreviewResources(filename, name, output)
		},
	}
	command.Flags().StringVarP(&name, "name", "n", "", "Name of the Application to preview")
	command.Flags().StringVarP(&output, "output", "o", "name", "Output format. One of: name|json|yaml")
	return command
}
