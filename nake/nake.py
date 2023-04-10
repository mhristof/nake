#! /usr/bin/env python3
# -*- coding: utf-8 -*-
# vim:fenc=utf-8
#

import argparse
import logging
import os

try:
    import precommit
    import make
except ImportError:
    import sys

    sys.path.append(os.path.dirname(__file__))
    import precommit
    import make

import hashlib


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("-C", default=os.getcwd(), help="Change to directory")
    parser.add_argument("-v", "--verbose", default=0, action="count", help="Verbose")

    args = parser.parse_args()

    level = logging.INFO

    if args.verbose > 0:
        level = logging.DEBUG

    logging.basicConfig(level=level)

    logging.debug("Changing to directory: %s", args.C)

    langs = languages(args.C)
    files = {
        ".pre-commit-config.yaml": precommit.render(langs),
        "Makefile": make.render(langs),
    }

    for filename, content in files.items():
        logging.debug("processing file: %s", filename)

        abs_file = os.path.join(args.C, filename)
        before_sha = None
        try:
            with open(filename, "rb") as f:
                logging.debug("Reading file: %s", abs_file)
                data = f.read()
                before_sha = hashlib.sha256(data).hexdigest()
        except FileNotFoundError:
            pass

        content_sha256 = hashlib.sha256(content.encode("utf-8")).hexdigest()

        if before_sha != content_sha256:
            logging.info("Updated %s", filename)

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
