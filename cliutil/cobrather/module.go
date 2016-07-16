package cobrather

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Module for cobra, used for cobra.Command
type Module struct {
	// The one-line usage message.
	Use string
	// The short description shown in the 'help' output.
	Short string
	// The long message shown in the 'help <this-command>' output.
	Long string
	// Examples of how to use the command
	Example string
	// RunE: Run but returns an error
	RunE func(cmd *cobra.Command, args []string) error
	// PostRunE: PostRun but returns an error
	PostRunE func(cmd *cobra.Command, args []string) error

	// extend cobra.Command fields
	GlobalFlags  []Flag
	Flags        []Flag
	Dependencies []*Module
	Commands     []*Module

	viper *viper.Viper
}

// MustNewCommand create cobra.Command from Module
func (t *Module) MustNewCommand() *cobra.Command {
	depModules := ListDeps(0, t)
	depSelfModules := append(depModules, t)

	command := &cobra.Command{
		Use:      t.Use,
		Short:    t.Short,
		Long:     t.Long,
		Example:  t.Example,
		PreRunE:  GenRunE(depModules...),
		RunE:     t.RunE,
		PostRunE: GenPostRunE(depSelfModules...),
	}

	for _, subcmd := range t.Commands {
		subcmd.viper = t.viper
		command.AddCommand(subcmd.MustNewCommand())
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
func (t *Module) MustNewRootCommand(v *viper.Viper) *cobra.Command {
	t.viper = v
	if t.viper == nil {
		t.viper = viper.New()
	}
	command := t.MustNewCommand()

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

// GenRunE generate RunE for all modules
func GenRunE(modules ...*Module) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) (err error) {
		walked := map[*Module]bool{}
		for _, module := range modules {
			if !walked[module] {
				walked[module] = true
				if module.RunE != nil {
					if err = module.RunE(cmd, args); err != nil {
						return
					}
				}
			}
		}
		return nil
	}
}

// GenPostRunE generate PostRunE for all modules, start from tail to head
func GenPostRunE(modules ...*Module) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) (err error) {
		walked := map[*Module]bool{}
		for i := len(modules) - 1; i >= 0; i-- {
			module := modules[i]
			if !walked[module] {
				walked[module] = true
				if module.PostRunE != nil {
					if err = module.PostRunE(cmd, args); err != nil {
						return
					}
				}
			}
		}
		return nil
	}
}
