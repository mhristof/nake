package cmd

import (
	"path/filepath"

	"github.com/mhristof/nake/bash"
	"github.com/mhristof/nake/docker"
	"github.com/spf13/cobra"
)

var (
	buildCmd = &cobra.Command{
		Use:   "build",
		Short: "Docker build",
		Run: func(cmd *cobra.Command, args []string) {
			dockerfile, _ := cmd.Flags().GetString("file")

			var docker = docker.Docker{
				Dockerfile: dockerfile,
				Pwd:        pwd,
			}

			name, _ := cmd.Flags().GetString("tag")

			var commands = []string{
				docker.Build(name),
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
	buildCmd.PersistentFlags().StringP("file", "f", "Dockerfile", "Dockerfile to use")
	abs, err := filepath.Abs(pwd)
	if err != nil {
		panic(err)
	}

	buildCmd.PersistentFlags().StringP("tag", "t", filepath.Base(abs), "The name of the image")

	rootCmd.AddCommand(buildCmd)
}
