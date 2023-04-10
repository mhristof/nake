#!/usr/bin/env bash
if [[ -z $CI ]] && [[ -z $AWS_PROFILE ]]; then
    echo "AWS_PROFILE not set"
    exit 1
fi

AWS_ACCOUNT_NAME=$(aws iam list-account-aliases --query "AccountAliases[0]" --output text)
# shellcheck disable=SC2001
ENV="$(sed 's/{{ company_name }}-//' <<<"${AWS_ACCOUNT_NAME}")"
VARS="-var-file=vars/$ENV.tfvars"

if [[ -z $AWS_REGION ]]; then
    AWS_REGION=$(aws configure get region)
fi

if [[ -f "vars/$ENV-$AWS_REGION.tfvars" ]]; then
    VARS="-var-file=vars/$ENV-$AWS_REGION.tfvars"
fi

export TF_DATA_DIR=".terraform-$ENV"

REMOTE=$(git config --get remote.origin.url | sed 's!https://.*@!!g' | sed 's!git@gitlab.com:\(.*\)!gitlab.com/\1!' | sed 's/.git$//')
export TF_CLI_ARGS_init="-backend-config=bucket=terraform-state-${AWS_ACCOUNT_NAME} -backend-config=key=$REMOTE/$AWS_REGION/terraform.tfstate"

export TF_CLI_ARGS_plan="$VARS"
export TF_CLI_ARGS_import="$VARS"
export TF_CLI_ARGS_console="$VARS"

export TF_VAR_env="$ENV"

GIT_ROOT=$(git rev-parse --show-toplevel)
TF_VAR_project="$(basename "$GIT_ROOT" | sed 's/-infra//g')"
export TF_VAR_project

if op --version &>/dev/null; then
    GITLAB_TOKEN="$(op --cache item get GITLAB_TOKEN --fields label=password)"
    export GITLAB_TOKEN
fi
