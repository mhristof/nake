#! /usr/bin/env python3
# -*- coding: utf-8 -*-
# vim:fenc=utf-8
#

import os
import logging
import jinja2

log = logging.getLogger(__name__)


def render(languages, defaults, features):
    default = "# vi: ft=bash\n"

    if "terraform-module" in features:
        languages.remove("terraform")

    for language in languages:
        try:
            with open(os.path.join(os.path.dirname(__file__), language + ".sh")) as f:
                template = f.read().strip().replace("#!/usr/bin/env bash", "")
                default += jinja2.Template(template).render(
                    {
                        "company_name": defaults["company"]["name"],
                        "gitlab_token": "terraform-gitlab-provider" in features,
                    }
                )
        except FileNotFoundError:
            log.debug("No config for language: %s", language)

    default = default.replace("#\n", "").replace("\n\n\n", "\n\n")

    if len(default.split("\n")) == 2:
        return None

    return default.strip() + "\n"
