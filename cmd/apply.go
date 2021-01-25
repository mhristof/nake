package cmd

import (
	"os"

	"github.com/mhristof/nake/bash"
	"github.com/mhristof/nake/log"
	"github.com/mhristof/nake/terraform"
	"github.com/spf13/cobra"
)

var (
	applyCmd = &cobra.Command{
		Use:   "apply",
		Short: "Terraform apply",
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

			commands := []string{
				tf.Init(),
				tf.Plan(),
				tf.Apply(),
			}
			for _, command := range commands {
				if command == "" {
					continue
				}
				bash.Eval(command)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(applyCmd)
}
