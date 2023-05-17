#! /usr/bin/env python3
# -*- coding: utf-8 -*-
# vim:fenc=utf-8
#


import logging
import os
import hcl2
import subprocess


log = logging.getLogger(__name__)


def render(pwd):
    lock = os.path.join(pwd, ".terraform.lock.hcl")

    if not os.path.exists(lock):
        return None

    with open(lock) as f:
        hcl = hcl2.load(f)

    tf_version = (
        subprocess.check_output(["terraform", "version"], stderr=subprocess.STDOUT)
        .decode("utf-8")
        .split(" ")[1]
        .split("\n")[0]
        .strip("v")
    )

    log.debug("Terraform version: '%s'", tf_version)

    providers = ""

    for provider in hcl["provider"]:
        name = list(provider.keys())[0]
        version = provider[name]["version"][0]

        source = name.replace("registry.terraform.io/", "")
        log.debug("Provider: %s %s", name, version)
        providers += f"""{os.path.basename(name)} = {{
        source = "{source}"
        version = ">= {version}"
        }}\n"""

    versions = f""" terraform {{
    required_version = ">= {tf_version}"
    required_providers {{
    {providers} }}
    }}"""

    formatted = subprocess.check_output(
        ["terraform", "fmt", "-"], input=versions.encode("utf-8")
    ).decode("utf-8")

    return formatted
