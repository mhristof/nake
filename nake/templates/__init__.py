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

    if "-infra" in remote:
        return terraform()

    return {}


def terraform():
    ret = {}

    template_dir = os.path.join(os.path.dirname(__file__), "terraform")

    log.debug("Reading templates from: %s", template_dir)

    for root, directory, files in os.walk(template_dir):
        log.debug("root: %s", root)

        for file in files:
            path = os.path.join(root, file)

            log.debug("Reading file: %s %s", root, file)

            dest = os.path.join(root, file).replace(template_dir + "/", "")
            with open(path, "r") as stream:
                ret[dest] = stream.read()

    return ret


def git_remote(directory):
    config = configparser.ConfigParser()
    config.read(os.path.join(directory, ".git", "config"))

    return config.get('remote "origin"', "url")
