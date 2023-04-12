#! /usr/bin/env python3
# -*- coding: utf-8 -*-
# vim:fenc=utf-8
#

import os
from . import render
import tempfile
import subprocess


def test_make_help():
    # new temp dir
    temp_dir = tempfile.mkdtemp()

    temp_file = os.path.join(temp_dir, "Makefile")

    with open(temp_file, "w") as stream:
        stream.write(render(["docker"], {"company": {"name": "test"}}))

    # os.chdir(temp_dir)
    help = subprocess.run(
        [f"make -C {temp_dir} help"],
        shell=True,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
    )

    assert help.returncode == 0
