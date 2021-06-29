package gnumake

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/mhristof/nake/log"
	"github.com/mhristof/nake/repo"
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
}

// RulesLib A library of make rules for known languages.
var RulesLib = map[string][]Rule{
	"help": []Rule{
		Rule{
			Help:    "Show this help",
			Targets: "help",
			Recipe:  "@grep '.*:.*##' Makefile | grep -v grep  | sort | sed 's/:.* ##/:/g' | column -t -s:",
			Phony:   true,
		},
	},
	"precommit": []Rule{
		Rule{
			Help:    "Install pre-commit checks",
			Targets: ".git/hooks/pre-commit",
			Recipe:  "pre-commit install",
		},
		Rule{
			Help:          "Run precommit checks",
			Targets:       "check",
			Recipe:        "pre-commit run --all",
			Prerequisites: ".git/hooks/pre-commit",
			Phony:         true,
		},
	},
	"Python": []Rule{
		Rule{
			Help:    "Run pep8 for the current directory",
			Targets: "pep8",
			Recipe:  "pycodestyle --ignore=E501",
			Phony:   true,
		},
	},
	"HCL": []Rule{
		Rule{
			Help:          "Force run 'terraform init'",
			Targets:       "init",
			Prerequisites: ".terraform",
			Phony:         true,
		},
		Rule{
			Targets: ".terraform",
			Recipe:  "terraform init",
		},
		Rule{
			Help:          "Runs 'terraform plan'",
			Targets:       "plan",
			Prerequisites: "terraform.tfplan",
			Phony:         true,
		},
		Rule{
			Help:          "Creates terraform.tfplan if required",
			Targets:       "terraform.tfplan",
			Prerequisites: "$(shell find ./ -name '*.tf')",
			Recipe:        "terraform plan -out $@",
		},
		Rule{
			Help:          "Run 'terraform apply'",
			Targets:       "apply",
			Prerequisites: "terraform.tfplan",
			Phony:         true,
		},
		Rule{
			Help:          "Run 'terraform apply' if required'",
			Targets:       "terraform.tfstate",
			Prerequisites: "terraform.tfplan",
			Recipe:        "terraform apply terraform.tfplan",
		},
		Rule{
			Help:    "Forcefully update terraform state",
			Targets: "force",
			Recipe:  "touch *.tf && make terraform.tfstate",
			Phony:   true,
		},
		Rule{
			Help:    "Run 'terraform destroy'",
			Targets: "destroy",
			Recipe:  "terraform destroy -auto-approve",
			Phony:   true,
		},
		Rule{
			Help:          "Clean the repository resources",
			Targets:       "clean",
			Prerequisites: "destroy",
			Recipe:        "rm -rf terraform.tf{state,plan} .terraform terraform.state.d",
			Phony:         true,
		},
	},
}

// Generate Print the Makefile based on the filetypes in the `dest` directory.
func Generate(dest string) {
	var rules Rules

	rules = append(rules, RulesLib["help"]...)

	for _, language := range repo.Languages(dest) {
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
		.DEFAULT_GOAL := help
		.ONESHELL:
	`),
	)
}

func (r Rules) pprint() {
	for _, rule := range r {
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
