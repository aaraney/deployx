package commands

import (
	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/command/completion"
	"github.com/docker/docker/api/types"
	"github.com/spf13/cobra"
)

// completeNames offers completion for swarm configs
func completeNames(dockerCli command.Cli) completion.ValidArgsFn {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		list, err := dockerCli.Client().ConfigList(cmd.Context(), types.ConfigListOptions{})
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		var names []string
		for _, config := range list {
			names = append(names, config.ID)
		}
		return names, cobra.ShellCompDirectiveNoFileComp
	}
}
