package main

import (
	"fmt"
	"os"

	"github.com/aaraney/deployx/commands"
	"github.com/aaraney/deployx/version"

	"github.com/docker/cli/cli"
	"github.com/docker/cli/cli-plugins/manager"
	"github.com/docker/cli/cli-plugins/plugin"
	"github.com/docker/cli/cli/command"
	cliflags "github.com/docker/cli/cli/flags"
)

func runStandalone(cmd *command.DockerCli) error {
	if err := cmd.Initialize(cliflags.NewClientOptions()); err != nil {
		return err
	}
	rootCmd := commands.NewDeployxCommand(cmd, false)
	return rootCmd.Execute()
}

func runPlugin(cmd *command.DockerCli) error {
	rootCmd := commands.NewDeployxCommand(cmd, true)
	return plugin.RunPlugin(cmd, rootCmd, manager.Metadata{
		SchemaVersion: "0.1.0",
		Vendor:        "aaraney", // TODO: not sure who to put here
		Version:       version.Version,
		URL:           "github.com/aaraney/deployx",
	})
}

func main() {
	cmd, err := command.NewDockerCli()
	if err != nil {
		fmt.Fprintln(cmd.Err(), err)
		os.Exit(1)
	}

	if plugin.RunningStandalone() {
		err = runStandalone(cmd)
	} else {
		err = runPlugin(cmd)
	}
	if err == nil {
		return
	}

	if sterr, ok := err.(cli.StatusError); ok {
		if sterr.Status != "" {
			fmt.Fprintln(cmd.Err(), sterr.Status)
		}
		// StatusError should only be used for errors, and all errors should
		// have a non-zero exit status, so never exit with 0
		if sterr.StatusCode == 0 {
			os.Exit(1)
		}
		os.Exit(sterr.StatusCode)
	}

	fmt.Fprintln(cmd.Err(), err)
	os.Exit(1)
}
