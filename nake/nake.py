#! /usr/bin/env python3
# -*- coding: utf-8 -*-
# vim:fenc=utf-8
#

import argparse
import logging
import os
import hashlib
import yaml
import importlib.metadata

import sys

sys.path.append(os.path.dirname(__file__))
import precommit
import make
import gitlabci
import envrc
import templates
import terraform

log = logging.getLogger(__name__)


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("-C", default=os.getcwd(), help="Change to directory")
    parser.add_argument("-v", "--verbose", default=0, action="count", help="Verbose")
    parser.add_argument(
        "--gitlab-token", default=os.getenv("GITLAB_TOKEN"), help="Gitlab token"
    )

    parser.add_argument(
        "--version",
        action="version",
        version="%(prog)s " + importlib.metadata.version("nake"),
    )

    parser.add_argument(
        "-c",
        "--config",
        default=os.path.join(os.getenv("XDG_CONFIG_HOME"), "nake.yml"),
        help="Config file",
    )

    args = parser.parse_args()

    level = logging.INFO

    if args.verbose > 0:
        level = logging.DEBUG

    logging.basicConfig(level=level)

    log.debug("Changing to directory: %s", args.C)

    conf = {}
    with open(args.config, "r") as stream:
        conf = yaml.safe_load(stream)

    save(templates.files(args.C), args)

    langs, features = languages(args.C)

    log.info("languages: %s", langs)

    envrcData = envrc.render(langs, conf, features)

    if envrcData is None:
        features |= {"no-envrc"}

    log.info("features: %s", features)

    files = {
        ".envrc": envrcData,
        ".pre-commit-config.yaml": precommit.render(langs),
        "Makefile": make.render(langs, conf),
        ".gitlab-ci.yml": gitlabci.render(
            args.C,
            args.gitlab_token,
            langs,
            conf[".gitlab-ci.yml"],
            features,
        ),
    }

    save(files, args)


def save(files, args):
    for filename, content in files.items():
        log.debug("processing file: %s", filename)

        abs_file = os.path.join(args.C, filename)

        if content is None and os.path.exists(abs_file):
            os.remove(abs_file)
            log.info("Removed %s", filename)

        if content is None:
            continue

        before_sha = None
        try:
            os.makedirs(os.path.dirname(abs_file), exist_ok=True)
            with open(abs_file, "rb") as f:
                log.debug("Reading file: %s", abs_file)
                data = f.read()
                before_sha = hashlib.sha256(data).hexdigest()
        except FileNotFoundError:
            pass

        content_sha256 = hashlib.sha256(content.encode("utf-8")).hexdigest()

        log.debug("sha before: %s", before_sha)
        log.debug("sha  after: %s", content_sha256)

        if before_sha != content_sha256:
            log.info("Updated %s", filename)

        with open(abs_file, "w") as stream:
            stream.write(content)


def file_as_bytes(file):
    with file:
        return file.read()


def languages(directory):
    ret = set()
    features = set()

    repo_name = os.path.basename(directory)

    if repo_name.startswith("terraform-"):
        features |= {"terraform-module"}

    for dirpath, dirnames, filenames in os.walk(directory):
        for filename in filenames:
            if ".terraform" in dirpath:
                continue

            if filename.endswith(".py"):
                ret |= {"python"}
            elif filename.endswith(".tf"):
                features |= has_gitlab_provider(os.path.join(dirpath, filename))
                ret |= {"terraform"}
            elif filename.endswith(".json"):
                ret |= {"json"}
            elif filename.startswith("Dockerfile"):
                log.debug("Found Dockerfile from %s/%s", dirpath, filename)
                ret |= {"docker"}
            elif filename.endswith(".go") and is_terratest_file(
                os.path.join(dirpath, filename)
            ):
                ret |= {"terratest"}

    return ret, features


def is_terratest_file(filename):
    log.debug("Checking if %s is terratest file", filename)
    with open(filename, "r") as stream:
        for line in stream:
            if "terraform.InitAndApply" in line:
                return True

    return False


def has_gitlab_provider(filename):
    with open(filename, "r") as stream:
        for line in stream:
            if "provider" in line and "gitlab" in line:
                return {"terraform-gitlab-provider"}

    return set()


if __name__ == "__main__":
    main()
