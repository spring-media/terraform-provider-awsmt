
venv:
	python -m venv $@

setup: venv
	set -eu && \
	. venv/bin/activate && \
	pip install --upgrade pip && \
	pip install --upgrade -r requirements.txt
.PHONY: setup

lint:
	set -eu && \
	. venv/bin/activate && \
	mkdocs build -f ../../mkdocs.yaml --strict
.PHONY: lint


clean:
	rm -rf venv
.PHONY: clean