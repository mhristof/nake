#! /usr/bin/env python3
# -*- coding: utf-8 -*-
# vim:fenc=utf-8
#

from . import *


def test_get_aws_region():
    assert get_aws_region("prod-ap-northeast-1") == "ap-northeast-1"
    assert get_aws_region("prod") == None
