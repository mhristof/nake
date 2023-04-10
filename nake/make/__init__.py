#! /usr/bin/env python3
# -*- coding: utf-8 -*-
# vim:fenc=utf-8
#

import logging
from textwrap import dedent
import os

log = logging.getLogger(__name__)


def render(languages):
    makefile = (
        dedent(
            """
        MAKEFLAGS += --warn-undefined-variables --jobs=$(shell nproc)
        SHELL := /bin/bash
        .SHELLFLAGS := -eu -o pipefail -c
        .DEFAULT_GOAL := build
        .ONESHELL:
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
        except FileNotFoundError:
            log.debug("No config for language: %s", language)

            continue

        makefile += "\n" + data.strip()
        added = True

    if not added:
        return None

    return makefile
