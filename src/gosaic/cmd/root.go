package cmd

import (
	"fmt"
	"gosaic/environment"
	"os"
	"path"
	"runtime"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// global flags
	dbPath  string
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
	home, err := homedir.Dir()
	if err != nil {
		fmt.Printf("Unable to get user home directory: %s\n", err.Error())
		os.Exit(1)
	}

	defaultDb := path.Join(home, ".gosaic.sqlite3")

	addGlobalStrFlag(&dbPath, "db", "", defaultDb, "Path to project database")
	addGlobalIntFlag(&workers, "workers", "", runtime.NumCPU(), "Number of workers to use")

	cobra.OnInitialize(setEnv)
}

func setEnv() {
	var err error
	Env, err = environment.GetProdEnv(
		viper.GetString("db"),
		viper.GetInt("workers"),
	)
	if err != nil {
		fmt.Printf("Unable to create environment: %s\n", err.Error())
		os.Exit(1)
	}
}

func addGlobalStrFlag(myVar *string, longName, shortName, defVal, desc string) {
	RootCmd.PersistentFlags().StringVarP(myVar, longName, shortName, defVal, desc)
	bindGlobalFlags(longName)
}

func addGlobalIntFlag(myVar *int, longName, shortName string, defVal int, desc string) {
	RootCmd.PersistentFlags().IntVarP(myVar, longName, shortName, defVal, desc)
	bindGlobalFlags(longName)
}

func addGlobalBoolFlag(myVar *bool, longName, shortName string, defVal bool, desc string) {
	RootCmd.PersistentFlags().BoolVarP(myVar, longName, shortName, defVal, desc)
	bindGlobalFlags(longName)
}

func bindGlobalFlags(flags ...string) {
	for _, flag := range flags {
		viper.BindPFlag(flag, RootCmd.PersistentFlags().Lookup(flag))
	}
}

func addLocalStrFlag(myVar *string, longName, shortName, defVal, desc string, cmds ...*cobra.Command) {
	for _, cmd := range cmds {
		cmd.Flags().StringVarP(myVar, longName, shortName, defVal, desc)
		bindLocalFlags(cmd, longName)
	}
}

func addLocalIntFlag(myVar *int, longName, shortName string, defVal int, desc string, cmds ...*cobra.Command) {
	for _, cmd := range cmds {
		cmd.Flags().IntVarP(myVar, longName, shortName, defVal, desc)
		bindLocalFlags(cmd, longName)
	}
}

func addLocalBoolFlag(myVar *bool, longName, shortName string, defVal bool, desc string, cmds ...*cobra.Command) {
	for _, cmd := range cmds {
		cmd.Flags().BoolVarP(myVar, longName, shortName, defVal, desc)
		bindLocalFlags(cmd, longName)
	}
}

func bindLocalFlags(cmd *cobra.Command, flags ...string) {
	for _, flag := range flags {
		viper.BindPFlag(flag, cmd.Flags().Lookup(flag))
	}
}
