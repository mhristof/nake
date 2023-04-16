#! /usr/bin/env python3
# -*- coding: utf-8 -*-
# vim:fenc=utf-8
#

import os
import configparser
import json
import logging


log = logging.getLogger(__name__)


def files(directory):
    remote = git_remote(directory)
    log.debug("Remote: %s", remote)

    files = render("default", directory)

    if "-infra" in remote:
        files.update(render("terraform", directory))

    if "docker-" in remote:
        files.update(render("docker", directory))

    return files


def render(folder, directory=""):
    ret = {}

    template_dir = os.path.join(os.path.dirname(__file__), folder)

    log.debug("Reading templates from: %s", template_dir)

    for root, _, files in os.walk(template_dir):
        log.debug("root: %s", root)

        for file in files:
            path = os.path.join(root, file)
            dest = os.path.join(root, file).replace(template_dir + "/", "")

            if os.path.exists(os.path.join(directory, dest)):
                log.debug("Skipping file: %s", dest)

                continue

            with open(path, "r") as stream:
                ret[dest] = stream.read()

    return ret


def git_remote(directory):
    config = configparser.ConfigParser(strict=False)
    config.read(os.path.join(directory, ".git", "config"))

    return config.get('remote "origin"', "url")
