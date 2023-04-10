#! /usr/bin/env python3
# -*- coding: utf-8 -*-
# vim:fenc=utf-8
#

import logging
import yaml
import os

log = logging.getLogger(__name__)


def render(languages):
    default_path = os.path.join(os.path.dirname(__file__), "default.yml")
    log.debug("Loading default config from: %s", default_path)

    default = None
    with open(default_path, "r") as stream:
        default = yaml.safe_load(stream)

    repos = default["repos"]

    for language in languages:
        try:
            with open(
                os.path.join(os.path.dirname(__file__), language + ".yml"), "r"
            ) as stream:
                repos += yaml.safe_load(stream)["repos"]
        except FileNotFoundError:
            log.debug("No config for language: %s", language)

    return yaml.dump({"repos": repos}, default_flow_style=False)
