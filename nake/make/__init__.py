#! /usr/bin/env python3
# -*- coding: utf-8 -*-
# vim:fenc=utf-8
#

import logging
from textwrap import dedent
import jinja2
import os
import re

log = logging.getLogger(__name__)


def render(languages, defaults):
    makefile = (
        dedent(
            """
            MAKEFLAGS += --warn-undefined-variables --jobs=$(shell nproc)
            SHELL := /bin/bash
            .SHELLFLAGS := -eu -o pipefail -c
            .DEFAULT_GOAL := build
            .ONESHELL:

            help:           ## Show this help.
                @grep '.*:.*##' Makefile | grep -v grep  | sort | sed 's/:.* ##/:/g' | column -t -s:
            """
        ).strip()
        + "\n"
    )

    added = False

    for language in languages:
        data = None
        try:
            with open(
                os.path.join(os.path.dirname(__file__), language + ".mk"), "r"
            ) as stream:
                data = stream.read()
                data = jinja2.Template(data).render(defaults)
        except FileNotFoundError:
            log.debug("No config for language: %s", language)

            continue

        makefile += "\n" + data.strip()
        added = True

    if not added:
        return None

    lines = []

    for line in makefile.splitlines():
        lines += [re.sub(r"^    ", "	", line)]

    return "\n".join(lines)
