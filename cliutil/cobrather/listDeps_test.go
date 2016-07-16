package cobrather_test

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tsaikd/KDGoLib/cliutil/cobrather"
)

func ExampleListDeps() {
	createModule := func(name string) *cobrather.Module {
		return &cobrather.Module{
			RunE: func(cmd *cobra.Command, args []string) error {
				fmt.Println(name)
				return nil
			},
		}
	}

	modCommonDep := createModule("modCommonDep")
	modCommon := createModule("modCommon")
	modCommon.Dependencies = []*cobrather.Module{modCommonDep}

	modCmdCmdDep := createModule("modCmdCmdDep")
	modCmdCmdDep.Dependencies = []*cobrather.Module{modCommon}
	modCmdCmd := createModule("modCmdCmd")
	modCmdCmd.Dependencies = []*cobrather.Module{modCmdCmdDep}

	modCmdDep := createModule("modCmdDep")
	modCmdDep.Dependencies = []*cobrather.Module{modCommon}
	modCmd := createModule("modCmd")
	modCmd.Dependencies = []*cobrather.Module{modCmdDep}
	modCmd.Commands = []*cobrather.Module{modCmdCmd}

	modRootDep := createModule("modRootDep")
	modRootDep.Dependencies = []*cobrather.Module{modCommon}
	modRoot := createModule("modRoot")
	modRoot.Dependencies = []*cobrather.Module{modRootDep}
	modRoot.Commands = []*cobrather.Module{modCmd}

	fmt.Println("only list dependencies of modRoot")
	for _, module := range cobrather.ListDeps(0, modRoot) {
		if err := module.RunE(nil, []string{}); err != nil {
			fmt.Println(err)
		}
	}
	fmt.Println()

	fmt.Println("list all dependencies of modRoot, include commands")
	for _, module := range cobrather.ListDeps(cobrather.OIncludeCommand, modRoot) {
		if err := module.RunE(nil, []string{}); err != nil {
			fmt.Println(err)
		}
	}
	fmt.Println()

	fmt.Println("list all dependencies of modRoot, include dependencies in commands, except commands")
	for _, module := range cobrather.ListDeps(cobrather.OIncludeDepInCommand, modRoot) {
		if err := module.RunE(nil, []string{}); err != nil {
			fmt.Println(err)
		}
	}
	fmt.Println()

	// Output:
	// only list dependencies of modRoot
	// modCommonDep
	// modCommon
	// modRootDep
	//
	// list all dependencies of modRoot, include commands
	// modCommonDep
	// modCommon
	// modRootDep
	// modCmdDep
	// modCmdCmdDep
	// modCmdCmd
	// modCmd
	//
	// list all dependencies of modRoot, include dependencies in commands, except commands
	// modCommonDep
	// modCommon
	// modRootDep
	// modCmdDep
	// modCmdCmdDep
}
