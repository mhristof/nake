#! /usr/bin/env python3
# -*- coding: utf-8 -*-
# vim:fenc=utf-8
#


import logging
import requests
import json
import yaml
import os
from functools import cmp_to_key

log = logging.getLogger(__name__)


def render(cwd, token, languages):
    config = {}
    try:
        with open(os.path.join(cwd, ".gitlab-ci.yml"), "r") as stream:
            config = yaml.safe_load(stream)
    except FileNotFoundError:
        pass

    if config is None:
        config = {
            "stages": [
                "lint",
                "plan",
                "apply",
                "release",
            ]
        }

    for language in languages:
        if language == "terraform":
            config = terraform(cwd, config)

    config["release"] = {
        **config.get("release", {}),
        **{
            "stage": "release",
            "script": ["npx semantic-release@20.1.1"],
            "variables": {"GITLAB_TOKEN": "$SEMANTIC_RELEASE_TOKEN"},
            "rules": [{"if": "$CI_COMMIT_BRANCH == $CI_DEFAULT_BRANCH"}],
        },
    }

    validate(token, config)

    return "---\n" + yaml.dump(
        config, Dumper=MyDumper, sort_keys=False, default_flow_style=False
    )


def terraform(cwd, config):
    config[".terraform"] = {
        **config.get(".terraform", {}),
        **{
            "image": {
                "name": "hashicorp/terraform:latest",
                "entrypoint": ["/bin/sh", "-c"],
            },
            "before_script": [
                "source .envrc",
                "set | grep '^TF'",
                "terraform init",
            ],
        },
    }

    config["fmt"] = {
        **config.get("fmt", {}),
        **{
            "stage": "lint",
            "script": "terraform fmt -check=true -recursive",
        },
    }

    config["yor"] = {
        **config.get("yor", {}),
        **{
            "image": {
                "name": "bridgecrew/yor:latest",
                "entrypoint": ["/bin/sh", "-c"],
            },
            "stage": "lint",
            "script": [
                "yor tag --tag-local-modules=true --directory .",
                "git diff --exit-code",
            ],
        },
    }

    tfvars = os.listdir(os.path.join(cwd, "vars"))

    if len(tfvars) == 0:
        return config

    sorted_tfvars = sorted(tfvars, key=lambda x: 0 if x.startswith("test.") else 1)

    plan_job = ""

    for tfvar in sorted_tfvars:
        name = os.path.splitext(tfvar)[0]

        log.debug("Found tfvar: %s", name)

        plan = yaml.safe_load(
            f"""
plan:
    resource_group: {name}
    environment:
      name: {name}
apply:
  extends: .terraform
  script:
    - source .envrc
    - terraform apply terraform.tfplan
  resource_group: {name}
  environment:
    name: {name}
  when: manual
  needs: ["{name}-plan"]
"""
        )

        if plan_job == "":
            plan_job = name
            plan["plan"]["extends"] = ".terraform"
            plan["plan"]["stage"] = "plan"
            plan["plan"]["script"] = [
                "terraform plan -out terraform.tfplan",
                """terraform show --json terraform.tfplan | jq -r '([.resource_changes[]?.change.actions?]|flatten)|{"create":(map(select(.=="create"))|length),"update":(map(select(.=="update"))|length),"delete":(map(select(.=="delete"))|length)}' > report.json""",
            ]
            plan["plan"]["artifacts"] = {
                "reports": {"terraform": ["report.json"]},
                "paths": ["terraform.tfplan"],
                "expire_in": "1 hour",
                "when": "always",
            }
            plan["apply"]["stage"] = "apply"
        else:
            plan["plan"]["extends"] = f"{plan_job}-plan"
            plan["apply"]["extends"] = f"{plan_job}-apply"
            plan["apply"]["needs"] += [f"{plan_job}-apply"]

        config = {
            **config,
            **{
                f"{name}-plan": plan["plan"],
                f"{name}-apply": plan["apply"],
            },
        }

        log.debug("Plan: %s", plan)

    return config


def validate(token, config):
    if token == "":
        log.warning("No GitLab token provided, skipping validation")

        return

    req = requests.post(
        "https://gitlab.com/api/v4/ci/lint?include_merged_yaml=true",
        json={"content": yaml.dump(config)},
        headers={
            "Content-Type": "application/json",
            "PRIVATE-TOKEN": "glpat-uFnbRGGHZRbb_N_x1z1i",
        },
    )

    if req.status_code != 200:
        raise Exception(f"Failed to validate config with status code {req.status_code}")

    if req.json()["status"] != "valid":
        raise Exception(f"Failed to validate config with error {req.json()['errors']}")


class MyDumper(yaml.SafeDumper):
    # HACK: insert blank lines between top-level objects
    # inspired by https://stackoverflow.com/a/44284819/3786245
    def write_line_break(self, data=None):
        super().write_line_break(data)

        if len(self.indents) == 1:
            super().write_line_break()
