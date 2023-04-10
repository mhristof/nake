#! /usr/bin/env python3
# -*- coding: utf-8 -*-
# vim:fenc=utf-8
#

import logging
from textwrap import dedent
import os


def render(languages):
    default = (
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

    for language in languages:
        data = None
        try:
            with open(
                os.path.join(os.path.dirname(__file__), language + ".mk"), "r"
            ) as stream:
                data = stream.read()
        except FileNotFoundError:
            logging.debug("make: No config for language: %s", language)

            continue

        default += "\n" + data.strip()

    logging.debug("default: %s", default)

    return default
