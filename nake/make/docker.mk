.PHONY: build
build: .build ## Build the docker image.

CI_JOB_STARTED_AT ?= $(shell date --iso-8601=seconds)
CI_COMMIT_SHA ?= $(shell git rev-parse HEAD)
CI_JOB_URL ?= $(USER)@$(shell hostname)
CI_PROJECT_NAME ?= $(shell basename $(shell pwd))
PROJECT_NAME ?= $(shell sed 's/docker-//' <<< $(CI_PROJECT_NAME))
CI_PROJECT_URL ?= $(shell git config --get remote.origin.url)
AWS_REGION ?= $(shell yq .variables.AWS_REGION .gitlab-ci.yml)
ECR ?= $(shell yq .variables.AWS_ACCOUNT_ID .gitlab-ci.yml).dkr.ecr.$(AWS_REGION).amazonaws.com
BASE_VERSION := $(shell grep FROM Dockerfile | tail -1  | awk '{print $$2}')

.build: Dockerfile $(shell grep COPY Dockerfile | sed 's/--from=\w*//' | cut -d ' ' -f2)
    docker build --progress=tty \
       --label=org.opencontainers.image.base.version=$(BASE_VERSION) \
       --label=org.opencontainers.image.created=$(CI_JOB_STARTED_AT) \
       --label=org.opencontainers.image.revision=$(CI_COMMIT_SHA) \
       --label=org.opencontainers.image.source=$(CI_JOB_URL) \
       --label=org.opencontainers.image.title=$(CI_PROJECT_NAME) \
       --label=org.opencontainers.image.url=$(CI_PROJECT_URL) \
       --label=org.opencontainers.image.vendor={{ company.name }} \
       --label=org.opencontainers.image.version=$(CI_COMMIT_SHA) \
       -t $(ECR)/$(PROJECT_NAME) .
    touch .build

.PHONY: run
run: ## Run the docker image.
    docker run --rm \
       --entrypoint /bin/bash \
       -v $(PWD):/work \
       -w /work \
       -it $(ECR)/$(PROJECT_NAME)
