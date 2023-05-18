#! /usr/bin/env python3
# -*- coding: utf-8 -*-
# vim:fenc=utf-8
#

from . import *


def test_get_aws_region():
    assert get_aws_region("prod-ap-northeast-1") == "ap-northeast-1"
    assert get_aws_region("prod") == None


def test_terraform():
    assert terraform({}, [])[0][".terraform"]["before_script"][0] == "source .envrc"
    assert (
        terraform({}, ["no-envrc"])[0][".terraform"]["before_script"][0]
        != "source .envrc"
    )


def test_terraform_variables():
    assert terraform_varfiles({}, ["test.tfvars"], [])[1] == ["plan", "apply"]

    assert (
        terraform_varfiles({}, ["prod-ap-northeast-1.tfvars"], [])[0][
            "prod-ap-northeast-1-apply"
        ]["variables"]["AWS_REGION"]
        == "ap-northeast-1"
    )

    assert (
        terraform_varfiles({}, ["test.tfvars", "prod-ap-northeast-1.tfvars"], [])[0][
            "prod-ap-northeast-1-apply"
        ]["variables"]["AWS_REGION"]
        == "ap-northeast-1"
    )

    assert (
        terraform_varfiles({}, ["test.tfvars", "prod.tfvars"], [])[0]["prod-plan"][
            "extends"
        ]
        == "test-plan"
    )

    assert (
        terraform_varfiles({}, ["test.tfvars"], [])[0]["test-apply"].get("when", None)
        == None
    )

    assert (
        terraform_varfiles({}, ["prod.tfvars"], [])[0]["prod-apply"].get("when", None)
        == "manual"
    )


def test_stages_compare():
    cases = {
        "simple job": (["stage", "script", "image"], ["stage", "image", "script"]),
    }

    for name, case in cases.items():
        case[0].sort(key=functools.cmp_to_key(stages_compare))

        assert case[0] == case[1]
