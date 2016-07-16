package cobrather_test

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tsaikd/KDGoLib/cliutil/cobrather"
)

func ExampleModule() {
	modCommon := &cobrather.Module{
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("modCommon Run")
			return nil
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("modCommon PostRun")
			return nil
		},
	}
	modForCmd := &cobrather.Module{
		Dependencies: []*cobrather.Module{modCommon},
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("modForCmd Run")
			return nil
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("modForCmd PostRun")
			return nil
		},
	}
	modForRoot := &cobrather.Module{
		Dependencies: []*cobrather.Module{modCommon},
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("modForRoot Run")
			return nil
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("modForRoot PostRun")
			return nil
		},
	}
	cmd := &cobrather.Module{
		Use:          "cmd",
		Dependencies: []*cobrather.Module{modForCmd},
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("cmd Run")
			return nil
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("cmd PostRun")
			return nil
		},
	}
	root := &cobrather.Module{
		Use:          "root",
		Dependencies: []*cobrather.Module{modForRoot},
		Commands:     []*cobrather.Module{cmd, cobrather.VersionModule},
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("root Run")
			return nil
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("root PostRun")
			return nil
		},
	}

	rootCommand := root.MustNewRootCommand(nil)
	if err := rootCommand.Execute(); err != nil {
		fmt.Println(err)
	}

	// Output:
	// modCommon Run
	// modForRoot Run
	// root Run
	// root PostRun
	// modForRoot PostRun
	// modCommon PostRun
}

func ExampleModule_cmd() {
	Viper := viper.New()
	flagCmd := &cobrather.StringFlag{
		Name:    "flagcmd",
		Default: "default flag cmd string",
		Usage:   "flag cmd string type",
	}
	modCommon := &cobrather.Module{
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("modCommon Run")
			return nil
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("modCommon PostRun")
			return nil
		},
	}
	modForCmd := &cobrather.Module{
		Dependencies: []*cobrather.Module{modCommon},
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("modForCmd Run")
			return nil
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("modForCmd PostRun")
			return nil
		},
	}
	modForRoot := &cobrather.Module{
		Dependencies: []*cobrather.Module{modCommon},
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("modForRoot Run")
			return nil
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("modForRoot PostRun")
			return nil
		},
	}
	cmd := &cobrather.Module{
		Use:          "cmd",
		Dependencies: []*cobrather.Module{modForCmd},
		GlobalFlags: []cobrather.Flag{
			flagCmd,
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("cmd Run", flagCmd.String())
			return nil
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("cmd PostRun")
			return nil
		},
	}
	root := &cobrather.Module{
		Use:          "root",
		Dependencies: []*cobrather.Module{modForRoot},
		Commands:     []*cobrather.Module{cmd, cobrather.VersionModule},
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("root Run")
			return nil
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("root PostRun")
			return nil
		},
	}

	rootCommand := root.MustNewRootCommand(Viper)
	rootCommand.SetArgs([]string{"cmd", "--flagcmd", "replace flag cmd"})
	if err := rootCommand.Execute(); err != nil {
		fmt.Println(err)
	}

	// Output:
	// modCommon Run
	// modForCmd Run
	// cmd Run replace flag cmd
	// cmd PostRun
	// modForCmd PostRun
	// modCommon PostRun
}

func ExampleModule_flag() {
	Viper := viper.New()
	flagFromArg := &cobrather.StringFlag{
		Name:    "fromarg",
		Default: "default flag arg string",
		Usage:   "flag string from arg",
	}
	flagFromEnv := &cobrather.StringFlag{
		Name:    "fromenv",
		Default: "default flag env string",
		Usage:   "flag string from env",
		EnvVar:  "FROMENV",
	}
	modCommon := &cobrather.Module{
		GlobalFlags: []cobrather.Flag{
			flagFromArg,
			flagFromEnv,
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("modCommon Run")
			return nil
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("modCommon PostRun")
			return nil
		},
	}
	modForCmd := &cobrather.Module{
		Dependencies: []*cobrather.Module{modCommon},
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("modForCmd Run")
			return nil
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("modForCmd PostRun")
			return nil
		},
	}
	modForRoot := &cobrather.Module{
		Dependencies: []*cobrather.Module{modCommon},
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("modForRoot Run")
			return nil
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("modForRoot PostRun")
			return nil
		},
	}
	cmd := &cobrather.Module{
		Use:          "cmd",
		Dependencies: []*cobrather.Module{modForCmd},
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("cmd Run")
			return nil
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("cmd PostRun")
			return nil
		},
	}
	root := &cobrather.Module{
		Use:          "root",
		Dependencies: []*cobrather.Module{modForRoot},
		Commands:     []*cobrather.Module{cmd, cobrather.VersionModule},
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("root Run")
			return nil
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("root PostRun")
			fmt.Println("from arg:", flagFromArg.String())
			fmt.Println("from env:", flagFromEnv.String())
			return nil
		},
	}

	os.Setenv("FROMENV", "test string from env")
	rootCommand := root.MustNewRootCommand(Viper)
	rootCommand.SetArgs([]string{"--fromarg", "test string from arg"})
	if err := rootCommand.Execute(); err != nil {
		fmt.Println(err)
	}

	// Output:
	// modCommon Run
	// modForRoot Run
	// root Run
	// root PostRun
	// from arg: test string from arg
	// from env: test string from env
	// modForRoot PostRun
	// modCommon PostRun
}
