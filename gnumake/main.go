package gnumake

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/MakeNowJust/heredoc"
	"github.com/mhristof/nake/repo"
	log "github.com/sirupsen/logrus"
)

// Rules A set of Makefile rules.
type Rules []Rule

// Rule Struct that holds all required info for a makefile rule.
type Rule struct {
	Targets       string
	Prerequisites string
	Recipe        string
	Phony         bool
	Help          string
	Variable      string
	Default       string
}

// RulesLib A library of make rules for known languages.
var RulesLib = map[string][]Rule{
	"help": {
		{
			Help:    "Show this help",
			Targets: "help",
			Recipe:  "@grep '.*:.*##' Makefile | grep -v grep  | sort | sed 's/:.* ##/:/g' | column -t -s:",
			Phony:   true,
		},
	},
	"git": {
		{
			Variable: `GIT_REF := $(shell git describe --match="" --always --dirty=+)`,
		},
		{
			Variable: `GIT_TAG := $(shell git name-rev --tags --name-only $(GIT_REF))`,
		},
	},
	"Go": {
		{
			Variable: `PACKAGE := $(shell go list)`,
		},
		{
			Targets: "test",
			Help:    "Run go test",
			Recipe:  "go test -v ./...",
			Phony:   true,
		},
		{
			Targets: fmt.Sprintf("%s.%s", bin(), runtime.GOOS),
			Help:    "Build the application binary for current OS",
		},
		{
			Default: fmt.Sprintf("%s.%s", bin(), runtime.GOOS),
		},
		{
			Help:    fmt.Sprintf("Build the application binary for target OS, for example %s.linux", bin()),
			Targets: fmt.Sprintf("%s.%%", bin()),
			Recipe:  `GOOS=$* go build -o $@ -ldflags "-X $(PACKAGE)/version=$(GIT_TAG)+$(GIT_REF)" main.go`,
		},
		{
			Help:          "Install the binary",
			Targets:       "install",
			Phony:         true,
			Prerequisites: fmt.Sprintf("%s.%s", bin(), runtime.GOOS),
			Recipe:        fmt.Sprintf("cp $< ~/bin/%s", packageName()),
		},
	},
	"precommit": {
		{
			Help:    "Install pre-commit checks",
			Targets: ".git/hooks/pre-commit",
			Recipe:  "pre-commit install",
		},
		{
			Help:          "Run precommit checks",
			Targets:       "check",
			Recipe:        "pre-commit run --all",
			Prerequisites: ".git/hooks/pre-commit",
			Phony:         true,
		},
	},
	"Python": {
		{
			Help:    "Run pep8 for the current directory",
			Targets: "pep8",
			Recipe:  "pycodestyle --ignore=E501",
			Phony:   true,
		},
	},
	"Docker": {
		{
			Targets:       ".build",
			Prerequisites: "Dockerfile",
			Recipe: heredoc.Doc(`
				docker build -t $(shell basename $(PWD) .
				touch .build
			`),
		},
		{
			Targets:       "build",
			Prerequisites: ".build",
			Help:          "Build the image",
		},
		{
			Targets:       "run",
			Prerequisites: ".build",
			Help:          "Run the image",
			Recipe:        "docker run $(shell basename $(PWD))",
		},
		{
			Targets:       "bash",
			Prerequisites: ".build",
			Help:          "Drop a shell into the image",
			Recipe:        "docker run -it --command /bin/bash $(shell basename $(PWD))",
		},
	},
	"terragrunt": {
		{
			Default: "all",
		},
		{
			Help:    "Apply all",
			Targets: "all",
			Phony:   true,
			Recipe:  "terragrunt run-all -auto-approve",
		},
		{
			Help:    "Clean all",
			Targets: "clean",
			Phony:   true,
			Recipe:  "terragrunt run-all destroy -auto-approve",
		},
		{
			Help:    "Plan all",
			Targets: "plan",
			Phony:   true,
			Recipe:  "terragrunt run-all plan",
		},
		{
			Help:    "Apply single folder",
			Targets: "%",
			Phony:   true,
			Recipe:  "cd $@ && terragrunt run-all plan",
		},
	},
	"HCL": {
		{
			Default: "terraform.tfstate",
		},
		{
			Help:          "Force run 'terraform init'",
			Targets:       "init",
			Prerequisites: ".terraform",
			Phony:         true,
		},
		{
			Targets: ".terraform",
			Recipe:  "terraform init",
		},
		{
			Help:          "Runs 'terraform plan'",
			Targets:       "plan",
			Prerequisites: "terraform.tfplan",
			Phony:         true,
		},
		{
			Help:          "Creates terraform.tfplan if required",
			Targets:       "terraform.tfplan",
			Prerequisites: "$(shell find ./ -name '*.tf') .terraform",
			Recipe:        "terraform plan -out $@",
		},
		{
			Help:          "Run 'terraform apply'",
			Targets:       "apply",
			Prerequisites: "terraform.tfstate",
			Phony:         true,
		},
		{
			Help:          "Run 'terraform apply' if required'",
			Targets:       "terraform.tfstate",
			Prerequisites: "terraform.tfplan",
			Recipe:        "terraform apply terraform.tfplan",
		},
		{
			Help:    "Forcefully update terraform state",
			Targets: "force",
			Recipe:  "touch *.tf && make terraform.tfstate",
			Phony:   true,
		},
		{
			Help:    "Run 'terraform destroy'",
			Targets: "destroy",
			Recipe:  "terraform destroy -auto-approve",
			Phony:   true,
		},
		{
			Help:          "Clean the repository resources",
			Targets:       "clean",
			Prerequisites: "destroy",
			Recipe:        "rm -rf terraform.tf{state,plan} .terraform terraform.state.d",
			Phony:         true,
		},
	},
}

// Generate Print the Makefile based on the filetypes in the `dest` directory.
func Generate(dest string, ignore []string) {
	var rules Rules

	rules = append(rules, RulesLib["help"]...)

	for _, language := range repo.Languages(dest, ignore) {
		log.WithFields(log.Fields{
			"language": language,
		}).Debug("Found")

		rules = append(rules, RulesLib[language]...)
	}

	rules.printHeader()
	rules.pprint()
}

func (r Rules) printHeader() {
	fmt.Println(heredoc.Doc(`
		MAKEFLAGS += --warn-undefined-variables
		SHELL := /bin/bash
		ifeq ($(word 1,$(subst ., ,$(MAKE_VERSION))),4)
		.SHELLFLAGS := -eu -o pipefail -c
		endif
		.ONESHELL:
	`),
	)
}

func (r Rules) pprint() {
	for _, rule := range r {
		if rule.Default == "" {
			continue
		}

		fmt.Println(fmt.Sprintf(".DEFAULT_GOAL := %s", rule.Default))
	}

	newLine := false
	// print variables on top
	for _, rule := range r {
		if rule.Variable == "" {
			continue
		}

		fmt.Println(rule.Variable)
		newLine = true
	}

	if newLine {
		fmt.Println("")
	}

	for _, rule := range r {
		if rule.Default != "" || rule.Variable != "" {
			continue
		}

		if rule.Phony {
			fmt.Println(fmt.Sprintf(".PHONY: %s", rule.Targets))
		}
		fmt.Println(fmt.Sprintf("%s: %s ## %s", rule.Targets, rule.Prerequisites, rule.Help))

		if rule.Recipe != "" {
			fmt.Println(fmt.Sprintf("\t%s", rule.Recipe))
		}
		fmt.Println("")
	}
}

func packageName() string {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	abs, err := filepath.Abs(pwd)
	if err != nil {
		panic(err)
	}

	return path.Base(abs)
}

func bin() string {
	return filepath.Join("bin/", packageName())
}
