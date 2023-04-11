.DEFAULT_GOAL := install

VERSION := $(shell yq .tool.poetry.version --input-format toml pyproject.toml -oy | tr -d 'v')
NEXT_VERSION := $(shell semver -n  | tail -1 | awk '{print $$NF}' | tr -d 'v')
GIT_TAG := $(shell git describe --tags --abbrev=0)

.PHONY: install
install: dist/nake-$(VERSION).tar.gz
	pip install dist/nake-$(VERSION).tar.gz

dist/nake-$(VERSION).tar.gz: pyproject.toml $(shell find . -type f -name '*.py')
	poetry build

release:
	if ! git diff --exit-code &> /dev/null; then echo "Working directory is dirty"; exit 1; fi
	sed -i "s/$(VERSION)/$(NEXT_VERSION)/" pyproject.toml
	git commit pyproject.toml -m "ci: release $(NEXT_VERSION)"
