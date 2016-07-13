package cmd

import (
	"fmt"
	"gosaic/environment"
	"os"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// global flags
var (
	path    string
	workers int
)

var (
	RootCmd = &cobra.Command{
		Use:   "gosaic",
		Short: "Create image mosaics",
		Long:  "Create image mosaics",
	}
	Env environment.Environment
)

func init() {
	cobra.OnInitialize(setEnv)

	addGlobalFlag(&path, "path", "p", "", "Path to project directory")
	addGlobalIntFlag(&workers, "workers", "w", runtime.NumCPU(), "Number of workers to use")
}

func setEnv() {
	var err error
	Env, err = environment.GetProdEnv(
		viper.GetString("path"),
		viper.GetInt("workers"),
	)
	if err != nil {
		fmt.Printf("Unable to create environment: %s\n", err.Error())
		os.Exit(1)
	}
}

func addGlobalFlag(myVar *string, longName, shortName, defVal, desc string) {
	RootCmd.PersistentFlags().StringVarP(myVar, longName, shortName, defVal, desc)
	bindGlobalFlags(longName)
}

func addGlobalIntFlag(myVar *int, longName, shortName string, defVal int, desc string) {
	RootCmd.PersistentFlags().IntVarP(myVar, longName, shortName, defVal, desc)
	bindGlobalFlags(longName)
}

func bindGlobalFlags(flags ...string) {
	for _, flag := range flags {
		viper.BindPFlag(flag, RootCmd.PersistentFlags().Lookup(flag))
	}
}

func addLocalFlag(myVar *string, longName, shortName, defVal, desc string, cmds ...*cobra.Command) {
	for _, cmd := range cmds {
		cmd.Flags().StringVarP(myVar, longName, shortName, defVal, desc)
		bindLocalFlags(cmd, longName)
	}
}

func bindLocalFlags(cmd *cobra.Command, flags ...string) {
	for _, flag := range flags {
		viper.BindPFlag(flag, cmd.Flags().Lookup(flag))
	}
}
