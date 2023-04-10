#! /usr/bin/env python3
# -*- coding: utf-8 -*-
# vim:fenc=utf-8
#


import logging
import functools
import requests
import io
import json
import ruamel.yaml
import os
from functools import cmp_to_key

log = logging.getLogger(__name__)
yaml = ruamel.yaml.YAML()


def render(cwd, token, languages, defaults):
    config = {}
    try:
        with open(os.path.join(cwd, ".gitlab-ci.yml"), "r") as stream:
            config = yaml.load(stream)
    except FileNotFoundError:
        pass

    if config is None:
        config = yaml.load("---\nstages: []\n")

    for language in languages:
        stages = []

        if language == "terraform":
            config, stages = terraform(cwd, config, defaults)
        elif language == "docker":
            config, stages = docker(cwd, config, defaults)

        config["stages"] = list(set(config.get("stages", []) + stages))
        log.debug("Stages: %s", config["stages"])

    config["stages"] = list(set(config.get("stages", []) + ["release"]))
    config["release"] = {
        **config.get("release", {}),
        **{
            "stage": "release",
            "script": ["npx semantic-release@20.1.1"],
            "variables": {"GITLAB_TOKEN": "$SEMANTIC_RELEASE_TOKEN"},
            "rules": [{"if": "$CI_COMMIT_BRANCH == $CI_DEFAULT_BRANCH"}],
        },
    }

    config["stages"].sort(key=functools.cmp_to_key(stages_compare))

    validate(token, config)

    return "---\n" + yaml_to_string(config)


def stages_compare(a, b):
    weights = {
        "lint": 10,
        "plan": 20,
        "build": 30,
        "apply": 40,
        "test": 50,
        "push": 60,
        "release": 70,
    }

    log.debug(
        "Comparing %s(%d) and %s(%d): %d"
        % (
            a,
            weights.get(a, 0),
            b,
            weights.get(b, 0),
            weights.get(a, 0) - weights.get(b, 0),
        )
    )

    return weights.get(a, 0) - weights.get(b, 0)


def docker(cwd, config, defaults):
    config["variables"] = {
        **config.get("variables", {}),
        **defaults["variables"],
        **{
            "ECR": "${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com",
        },
    }

    config["image"] = config.get("image", "docker:latest")

    config[".docker-auth"] = {
        **config.get(".docker-auth", {}),
        **{
            "before_script": [
                "aws ecr get-login-password --region $AWS_REGION | docker login --username AWS --password-stdin $ECR",
            ],
        },
    }

    project_name = os.path.basename(cwd).replace("docker-", "")
    config["build"] = {
        **config.get("build", {}),
        **{
            "extends": [".docker-auth"],
            "stage": "build",
            "script": [
                "apt install -y make",
                "make build",
            ],
        },
    }

    config["push"] = {
        **config.get("push", {}),
        **{
            "extends": [".docker-auth"],
            "stage": "push",
            "script": [
                "apt install -y make",
                "make build",
                "docker tag $ECR/%s:$CI_COMMIT_SHORT_SHA $ECR/${project_name}:${CI_COMMIT_TAG/v/}"
                % project_name,
                "docker push $ECR/%s:${CI_COMMIT_TAG/v/}" % project_name,
            ],
            "rules": [{"if": "$CI_COMMIT_TAG"}],
        },
    }

    return config, ["build", "push"]


def terraform(cwd, config, defaults):
    stages = ["lint"]

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
        **defaults.get(".terraform", {}),
    }

    config["fmt"] = {
        **config.get("fmt", {}),
        **{
            "stage": "lint",
            "extends": [".terraform"],
            "script": [
                "terraform fmt -diff -check=true -recursive",
            ],
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
        **defaults.get("yor", {}),
    }

    tfvars = os.listdir(os.path.join(cwd, "vars"))

    if len(tfvars) == 0:
        return config, stages

    stages += ["plan", "apply"]

    sorted_tfvars = sorted(tfvars, key=lambda x: 0 if x.startswith("test.") else 1)

    plan_job = ""

    for tfvar in sorted_tfvars:
        name = os.path.splitext(tfvar)[0]

        log.debug("Found tfvar: %s", name)

        plan = yaml.load(
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

    return config, stages


def yaml_to_string(data):
    buf = io.BytesIO()
    yaml.indent(mapping=2, sequence=4, offset=2)
    yaml.preserve_quotes = True
    yaml.width = 4096
    yaml.dump(data, buf)

    new_lines = []

    for line in buf.getvalue().decode("utf-8").splitlines():
        if line == "":
            continue

        if not line.startswith(" "):
            new_lines.append("")
            new_lines.append(line)

            continue

        new_lines.append(line)

    return "\n".join(new_lines).strip() + "\n"


def validate(token, config):
    if token is None:
        log.warning("No GitLab token provided, skipping validation")

        return

    req = requests.post(
        "https://gitlab.com/api/v4/ci/lint?include_merged_yaml=true",
        json={"content": yaml_to_string(config)},
        headers={
            "Content-Type": "application/json",
            "PRIVATE-TOKEN": token,
        },
    )

    if req.status_code != 200:
        raise Exception(f"Failed to validate config with status code {req.status_code}")

    if req.json()["status"] != "valid":
        with open("gitlab-ci.yml", "w") as f:
            f.write(yaml_to_string(config))

        raise Exception(f"Failed to validate config with error {req.json()['errors']}")
