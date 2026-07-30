
venv:
	python -m venv $@

setup: venv
	pip install --upgrade pip && \
	pip install --upgrade -r requirements.txt
.PHONY: setup

lint:
	mkdocs build -f ../../mkdocs.yaml --strict
.PHONY: lint


clean:
	rm -rf venv
.PHONY: clean