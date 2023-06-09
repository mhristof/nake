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


def dockerfile(languages, config):
    print(config)

    C = config["C"]

    if "go" in languages and imports_lambda(os.path.join(C, "main.go")):
        return f"""FROM golang:{go_version(C)}-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main .

FROM public.ecr.aws/lambda/go:1.2023.05.28.19

COPY --from=builder /app/main /var/task/main
CMD [ "main" ]
"""


def go_version(directory):
    with open(os.path.join(directory, "go.mod"), "r") as stream:
        for line in stream.readlines():
            if "go " in line:
                return line.split(" ")[1].strip()

    raise Exception("Could not find go version")


def imports_lambda(file):
    with open(file, "r") as stream:
        for line in stream.readlines():
            if "github.com/aws/aws-lambda-go/lambda" in line:
                return True

    return False


def git_remote(directory):
    config = configparser.ConfigParser(strict=False)
    config.read(os.path.join(directory, ".git", "config"))

    return config.get('remote "origin"', "url")
