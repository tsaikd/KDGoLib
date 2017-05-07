package tcobrather

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tsaikd/KDGoLib/cliutil/cobrather"
)

func Example() {
	flagDepArg := &cobrather.StringFlag{
		Name:   "deparg",
		EnvVar: "DEPARG",
	}
	flagRootArg := &cobrather.StringFlag{
		Name:    "rootarg",
		Default: "default root arg",
		EnvVar:  "ROOTARG",
	}
	moduleRootDep1Dep := &cobrather.Module{
		Flags: []cobrather.Flag{
			flagDepArg,
		},
		RunE: func(ctx context.Context, cmd *cobra.Command, args []string) error {
			fmt.Println("root dep1 dep start", flagDepArg.String())
			return nil
		},
		PostRunE: func(ctx context.Context, cmd *cobra.Command, args []string) error {
			fmt.Println("root dep1 dep close")
			return nil
		},
	}
	moduleRootDep1 := &cobrather.Module{
		Dependencies: []*cobrather.Module{
			moduleRootDep1Dep,
		},
		RunE: func(ctx context.Context, cmd *cobra.Command, args []string) error {
			fmt.Println("root dep1 start")
			return nil
		},
		PostRunE: func(ctx context.Context, cmd *cobra.Command, args []string) error {
			fmt.Println("root dep1 close")
			return nil
		},
	}
	moduleRootDep2 := &cobrather.Module{
		RunE: func(ctx context.Context, cmd *cobra.Command, args []string) error {
			fmt.Println("root dep2 start")
			return nil
		},
		PostRunE: func(ctx context.Context, cmd *cobra.Command, args []string) error {
			fmt.Println("root dep2 close")
			return nil
		},
	}
	moduleRoot := &cobrather.Module{
		Dependencies: []*cobrather.Module{
			moduleRootDep1,
			moduleRootDep2,
		},
		Flags: []cobrather.Flag{
			flagRootArg,
		},
		RunE: func(ctx context.Context, cmd *cobra.Command, args []string) error {
			fmt.Println("root start", flagRootArg.String())
			return nil
		},
		PostRunE: func(ctx context.Context, cmd *cobra.Command, args []string) error {
			fmt.Println("root close")
			return nil
		},
	}

	if err := os.Setenv("DEPARG", "dep arg config"); err != nil {
		fmt.Println(err)
	}
	ctx := context.Background()
	testmod := NewTest(ctx, moduleRoot)
	if err := testmod.Setup(); err != nil {
		fmt.Println(err)
	}
	fmt.Println("do some testing")
	if err := testmod.Teardown(); err != nil {
		fmt.Println(err)
	}

	// Output:
	// root dep1 dep start dep arg config
	// root dep1 start
	// root dep2 start
	// root start default root arg
	// do some testing
	// root close
	// root dep2 close
	// root dep1 close
	// root dep1 dep close
}
