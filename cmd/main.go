package main

import cmd "github.com/touchardv/argocd-offline-cli/cmd/commands"

func main() {
	command := cmd.NewCommand()
	command.Execute()
}
