package commands

import (
	"errors"
	"os"

	"github.com/aaraney/deployx/commands/options"
	"github.com/aaraney/deployx/deploy"

	composego_cli "github.com/compose-spec/compose-go/cli"
	"github.com/compose-spec/compose-go/loader"
	composego "github.com/compose-spec/compose-go/types"
	"github.com/docker/cli/cli"
	"github.com/docker/cli/cli-plugins/plugin"
	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/command/stack/swarm"
	"github.com/spf13/cobra"
)

func NewDeployxCommand(dockerCli command.Cli, isPlugin bool) *cobra.Command {
	var opts options.Deploy

	cmd := &cobra.Command{
		Use:   "deployx [OPTIONS] STACK",
		Short: "Deploy a new stack or update an existing stack",
		Args:  cli.RequiresMaxArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				cmd.Usage()
				return nil
			}
			opts.Namespace = args[0]
			return runDeployxCommand(dockerCli, &opts)
		},
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return completeNames(dockerCli)(cmd, args, toComplete)
		},
	}
	if isPlugin {
		cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
			return plugin.PersistentPreRunE(cmd, args)
		}
	} else {
		// match plugin behavior for standalone mode
		// https://github.com/docker/cli/blob/6c9eb708fa6d17765d71965f90e1c59cea686ee9/cli-plugins/plugin/plugin.go#L117-L127
		cmd.SilenceUsage = true
		cmd.SilenceErrors = true
		cmd.TraverseChildren = true
		cmd.DisableFlagsInUseLine = true
		cli.DisableFlagsInUseLine(cmd)
	}

	flags := cmd.Flags()
	flags.StringSliceVarP(&opts.Composefiles, "compose-file", "c", []string{}, `Path to a Compose file, or "-" to read from stdin`)
	flags.SetAnnotation("compose-file", "version", []string{"2"})
	flags.StringSliceVar(&opts.Envfiles, "env-file", []string{}, `Path to an alternative env file, or "-" to read from stdin`)
	flags.BoolVar(&opts.NoInterpolate, "no-interpolate", false, "Don't perform environment variable interpolation")
	flags.BoolVar(&opts.SendRegistryAuth, "with-registry-auth", false, "Send registry authentication details to Swarm agents")
	flags.BoolVar(&opts.Prune, "prune", false, "Prune services that are no longer referenced")
	flags.SetAnnotation("prune", "version", []string{"1.27"})
	flags.StringVar(&opts.ResolveImage, "resolve-image", swarm.ResolveImageAlways,
		`Query the registry to resolve image digest and supported platforms ("`+swarm.ResolveImageAlways+`", "`+swarm.ResolveImageChanged+`", "`+swarm.ResolveImageNever+`")`)
	flags.SetAnnotation("resolve-image", "version", []string{"1.30"})
	return cmd
}

func getEnv(workingDir string, envFiles []string) (map[string]string, error) {
	env := environMap()
	new_env, err := composego_cli.GetEnvFromFile(env, workingDir, envFiles)
	if err != nil {
		return env, err
	}

	// TODO: right now intersection does not overwrite. consider if it should
	leftMerge(env, new_env)

	return env, nil
}

func runDeployxCommand(dockerCli command.Cli, opts *options.Deploy) error {
	if opts == nil {
		return errors.New("options.Deploy is nil")
	}

	if err := validateStackName(opts.Namespace); err != nil {
		return err
	}

	workingDir, err := os.Getwd()
	if err != nil {
		return err
	}

	env := make(map[string]string)
	if !opts.NoInterpolate {
		// NOTE: if no env files are provided, .env file is default
		env, err = getEnv(workingDir, opts.Envfiles)
		if err != nil {
			return err
		}
	}

	if len(opts.Composefiles) < 1 {
		return errors.New("Please specify a Compose file (with --compose-file)")
	}

	var compose_files []composego.ConfigFile
	for _, compose_file := range opts.Composefiles {
		compose_files = append(compose_files, composego.ConfigFile{Filename: compose_file})
	}

	config, err := loader.Load(composego.ConfigDetails{
		WorkingDir:  workingDir,
		ConfigFiles: compose_files,
		Environment: env,
	}, func(options *loader.Options) {
		options.SetProjectName(opts.Namespace, true)
	})
	if err != nil {
		return err
	}

	cfg := composego.Config{
		// Filename:   filename, // note sure if this is right since this could be merged compose files
		Name:       config.Name,
		Services:   config.Services,
		Networks:   config.Networks,
		Volumes:    config.Volumes,
		Secrets:    config.Secrets,
		Configs:    config.Configs,
		Extensions: config.Extensions,
	}
	return deploy.RunDeploy(dockerCli, *opts, &cfg)
}
