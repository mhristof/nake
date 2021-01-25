package cmd

import (
	"os"

	"github.com/mhristof/nake/bash"
	"github.com/mhristof/nake/log"
	"github.com/mhristof/nake/terraform"
	"github.com/spf13/cobra"
)

var (
	initCmd = &cobra.Command{
		Use:   "init",
		Short: "Terraform init",
		Run: func(cmd *cobra.Command, args []string) {
			dir, err := os.Getwd()
			if err != nil {
				log.Fatal(err)
			}

			var tf = terraform.Terraform{
				Pwd: dir,
			}

			if !tf.Available() {
				log.Fatal("Terraform not available for this folder")
			}

			command := tf.Init()

			if command != "" {
				bash.Eval(command)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(initCmd)
}
