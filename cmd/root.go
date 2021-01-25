package cmd

import (
	"fmt"
	"os"

	"github.com/mhristof/nake/log"
	"github.com/spf13/cobra"
)

var version = "devel"
var pwd = ""

var rootCmd = &cobra.Command{
	Use:     "nake",
	Short:   "Opinionated make to handle common DevOps patterns",
	Long:    `TODO: changeme`,
	Version: version,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		Verbose(cmd)
		pwd, _ = os.Getwd()
	},
}

// Verbose Increase verbosity
func Verbose(cmd *cobra.Command) {
	verbose, err := cmd.Flags().GetBool("verbose")
	if err != nil {
		log.Panic(err)
	}

	if verbose {
		log.SetLevel(log.DebugLevel)
	}
}
func init() {
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Increase verbosity")
	rootCmd.PersistentFlags().BoolP("dryrun", "n", false, "Dry run")
}

// Execute The main function for the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
