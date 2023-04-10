#! /usr/bin/env python3
# -*- coding: utf-8 -*-
# vim:fenc=utf-8
#

import argparse
import logging
import os
import hashlib
import yaml

import sys

sys.path.append(os.path.dirname(__file__))
import precommit
import make
import gitlabci


log = logging.getLogger(__name__)


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("-C", default=os.getcwd(), help="Change to directory")
    parser.add_argument("-v", "--verbose", default=0, action="count", help="Verbose")
    parser.add_argument(
        "--gitlab-token", default=os.getenv("GITLAB_TOKEN"), help="Gitlab token"
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

    langs = languages(args.C)
    files = {
        ".pre-commit-config.yaml": precommit.render(langs),
        "Makefile": make.render(langs),
        ".gitlab-ci.yml": gitlabci.render(
            args.C, args.gitlab_token, langs, conf[".gitlab-ci.yml"]
        ),
    }

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

    for dirpath, dirnames, filenames in os.walk(directory):
        for filename in filenames:
            if filename.endswith(".py"):
                ret |= {"python"}
            elif filename.endswith(".tf"):
                ret |= {"terraform"}
            elif filename.endswith(".json"):
                ret |= {"json"}
            elif filename.startswith("Dockerfile"):
                ret |= {"docker"}

    return ret


if __name__ == "__main__":
    main()
