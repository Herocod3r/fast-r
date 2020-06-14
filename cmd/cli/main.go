package main

import (
	"github.com/herocod3r/fast-r/cmd/cli/commands"
	"github.com/spf13/cobra"
)

const (
	cliName        = "Fast-R Speed Test CLI"
	cliDescription = "A simple command line client for speed testing."
)

var (
	rootCmd = &cobra.Command{
		Use:        cliName,
		Short:      cliDescription,
		SuggestFor: []string{"fast-r"},
	}
)

func main() {
	//Client Info
	//run
	rootCmd.AddCommand(
		commands.NewInfoCommand(),
		commands.NewRunCommand(),
	)
	cobra.EnablePrefixMatching = true
	if err := rootCmd.Execute(); err != nil {

	}
}
