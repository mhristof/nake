
.DEFAULT_GOAL := install

.PHONY: install
install: dist/nake-0.1.0.tar.gz
	pip install dist/nake-0.1.0.tar.gz

dist/nake-0.1.0.tar.gz: pyproject.toml $(shell find . -type f -name '*.py')
	poetry build
