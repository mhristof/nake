#! /usr/bin/env python3
# -*- coding: utf-8 -*-
# vim:fenc=utf-8
#


import logging
import functools
import requests
import io
import re
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
            config, stages = terraform(cwd, config)
        elif language == "docker":
            config, stages = docker(cwd, config)

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

    for k, v in defaults.items():
        for k2, v2 in config.items():
            if re.match(k, k2):
                config[k2] = {**v2, **v}

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


def docker(cwd, config):
    config["variables"] = {
        **config.get("variables", {}),
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


def terraform(cwd, config):
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
                "name": "bridgecrew/yor:0.1.170",
                "entrypoint": [
                    "/usr/bin/env",
                    "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
                ],
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
        return config, stages

    stages += ["plan", "apply"]

    sorted_tfvars = sorted(tfvars, key=lambda x: 0 if x.startswith("test.") else 1)

    plan_job = ""

    for tfvar in sorted_tfvars:
        name = os.path.splitext(tfvar)[0]
        region = get_aws_region(name)

        log.debug("Found tfvar: %s with region: %s", name, region)

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
            plan["plan"]["needs"] = []
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

            if region is not None:
                plan["plan"]["variables"] = {
                    **plan["plan"].get("variables", {}),
                    **{"AWS_REGION": region},
                }
                plan["apply"]["variables"] = {
                    **plan["apply"].get("variables", {}),
                    **{"AWS_REGION": region},
                }
                log.debug("Setting AWS_REGION to %s", region)

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
    yaml = ruamel.yaml.YAML()
    buf = io.BytesIO()
    yaml.indent(mapping=2, sequence=4, offset=2)
    yaml.preserve_quotes = True
    yaml.Representer = NonAliasingRTRepresenter
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


def is_aws_regions(region):
    all_regions = """\
        us-east-2
        us-east-1
        us-west-1
        us-west-2
        af-south-1
        ap-east-1
        ap-south-2
        ap-southeast-3
        ap-southeast-4
        ap-south-1
        ap-northeast-3
        ap-northeast-2
        ap-southeast-1
        ap-southeast-2
        ap-northeast-1
        ca-central-1
        eu-central-1
        eu-west-1
        eu-west-2
        eu-south-1
        eu-west-3
        eu-south-2
        eu-north-1
        eu-central-2
        me-south-1
        me-central-1
        sa-east-1
    """.strip().split()

    log.debug("Regions: %s", all_regions)

    ret = region in all_regions
    log.debug("region %s (is region: %s)", region, ret)
    return ret


def get_aws_region(name):
    region = None

    try:
        region = "-".join(name.split("-")[1:])
    except IndexError:
        log.debug("Failed to get region for %s", name)
        return None

    if region is None:
        log.debug("Failed to get region for %s", name)
        return None

    if not is_aws_regions(region):
        log.debug("Not a region: %s", region)
        return None

    return region


class NonAliasingRTRepresenter(ruamel.yaml.representer.RoundTripRepresenter):
    def ignore_aliases(self, data):
        return True
