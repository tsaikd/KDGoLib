package cobrather

import (
	"context"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Module for cobra, used for cobra.Command
type Module struct {
	// The one-line usage message.
	Use string
	// An array of aliases that can be used instead of the first word in Use.
	Aliases []string
	// The short description shown in the 'help' output.
	Short string
	// The long message shown in the 'help <this-command>' output.
	Long string
	// Examples of how to use the command
	Example string
	// RunE: Run but returns an error
	RunE func(ctx context.Context, cmd *cobra.Command, args []string) error
	// PostRunE: PostRun but returns an error
	PostRunE func(ctx context.Context, cmd *cobra.Command, args []string) error

	// extend cobra.Command fields
	GlobalFlags  []Flag
	Flags        []Flag
	Dependencies []*Module
	Commands     []*Module

	viper *viper.Viper
}

// MustNewCommand create cobra.Command from Module
func (t *Module) MustNewCommand(ctx context.Context) *cobra.Command {
	depModules := ListDeps(0, t)
	depSelfModules := append(depModules, t)

	command := &cobra.Command{
		Use:     t.Use,
		Aliases: t.Aliases,
		Short:   t.Short,
		Long:    t.Long,
		Example: t.Example,
		PreRunE: GenRunE(ctx, depModules...),
		RunE: func(cmd *cobra.Command, args []string) error {
			select {
			case <-ctx.Done():
				return nil
			default:
			}
			return t.RunE(ctx, cmd, args)
		},
		PostRunE: GenPostRunE(ctx, depSelfModules...),
	}

	for _, subcmd := range t.Commands {
		subcmd.viper = t.viper
		command.AddCommand(subcmd.MustNewCommand(ctx))
	}

	for _, module := range depSelfModules {
		for _, flag := range module.Flags {
			if err := flag.Bind(command.Flags(), t.viper); err != nil {
				panic(err)
			}
		}
	}

	return command
}

// MustNewRootCommand create cobra.Command from Module for root application
func (t *Module) MustNewRootCommand(ctx context.Context, vr *viper.Viper) *cobra.Command {
	t.viper = vr
	if t.viper == nil {
		t.viper = viper.New()
	}
	command := t.MustNewCommand(ctx)

	modules := append(ListDeps(OIncludeCommand, t), t)
	for _, module := range modules {
		for _, flag := range module.GlobalFlags {
			if err := flag.Bind(command.PersistentFlags(), t.viper); err != nil {
				panic(err)
			}
		}
	}

	return command
}

// MainRunOption option for Module.MustMainRun
type MainRunOption interface{}

// MainRunOptionContext config Context interface for running commands
type MainRunOptionContext context.Context

// MainRunOptionViper config Viper
type MainRunOptionViper *viper.Viper

// MainRunOptionSilenceUsage config SilenceUsage before run, default true
type MainRunOptionSilenceUsage bool

// MustMainRun used for main package to run main module
func (t *Module) MustMainRun(options ...MainRunOption) {
	var vr *viper.Viper
	ctx := context.Background()
	SilenceUsage := true

	for _, option := range options {
		switch opt := option.(type) {
		case MainRunOptionContext:
			ctx = opt
		case MainRunOptionViper:
			vr = opt
		case MainRunOptionSilenceUsage:
			SilenceUsage = bool(opt)
		}
	}

	command := t.MustNewRootCommand(ctx, vr)
	command.SilenceUsage = SilenceUsage
	if err := command.Execute(); err != nil {
		os.Exit(-1)
	}
}

// GenRunE generate RunE for all modules
func GenRunE(ctx context.Context, modules ...*Module) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) (err error) {
		walked := map[*Module]bool{}
		for _, module := range modules {
			if !walked[module] {
				walked[module] = true
				if module.RunE != nil {
					if err = module.RunE(ctx, cmd, args); err != nil {
						return
					}
				}
			}
		}
		return nil
	}
}

// GenPostRunE generate PostRunE for all modules, start from tail to head
func GenPostRunE(ctx context.Context, modules ...*Module) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) (err error) {
		walked := map[*Module]bool{}
		for i := len(modules) - 1; i >= 0; i-- {
			module := modules[i]
			if !walked[module] {
				walked[module] = true
				if module.PostRunE != nil {
					if err = module.PostRunE(ctx, cmd, args); err != nil {
						return
					}
				}
			}
		}
		return nil
	}
}
