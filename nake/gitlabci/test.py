#! /usr/bin/env python3
# -*- coding: utf-8 -*-
# vim:fenc=utf-8
#

from . import *


def test_get_aws_region():
    assert get_aws_region("prod-ap-northeast-1") == "ap-northeast-1"
    assert get_aws_region("prod") == None


def test_terraform_variables():
    assert terraform_varfiles({}, ["test.tfvars"])[1] == ["plan", "apply"]

    assert (
        terraform_varfiles({}, ["prod-ap-northeast-1.tfvars"])[0][
            "prod-ap-northeast-1-apply"
        ]["variables"]["AWS_REGION"]
        == "ap-northeast-1"
    )

    assert (
        terraform_varfiles({}, ["test.tfvars", "prod-ap-northeast-1.tfvars"])[0][
            "prod-ap-northeast-1-apply"
        ]["variables"]["AWS_REGION"]
        == "ap-northeast-1"
    )

    assert (
        terraform_varfiles({}, ["test.tfvars", "prod.tfvars"])[0]["prod-plan"][
            "extends"
        ]
        == "test-plan"
    )
