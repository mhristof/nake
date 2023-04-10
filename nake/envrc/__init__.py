#! /usr/bin/env python3
# -*- coding: utf-8 -*-
# vim:fenc=utf-8
#

import os
import logging
import jinja2

log = logging.getLogger(__name__)


def render(languages, defaults):
    default = "# vi: ft=bash\n"

    for language in languages:
        try:
            with open(os.path.join(os.path.dirname(__file__), language + ".sh")) as f:
                template = f.read().strip().replace("#!/usr/bin/env bash", "")
                default += jinja2.Template(template).render(
                    {
                        "company_name": defaults["company"]["name"],
                    }
                )
        except FileNotFoundError:
            log.debug("No config for language: %s", language)

    return default.strip() + "\n"
