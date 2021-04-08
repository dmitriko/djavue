# simple make file for django/vue  project
.PHONY: venv-update test build-backend migrate pepare-instance deploy-frontend deploy-backend

PIP = venv/bin/pip
PY = venv/bin/python
TAG := $(shell git describe)

venv:
	python3 -m venv venv

venv-update: venv
	$(PIP) install --upgrade -r requirements.txt
	$(PIP) install --upgrade -r requirements-dev.txt

test: venv
	$(PY) manage.py test

migrate:
	$(PY) manage.py makemigrations
	$(PY) manage.py migrate

build-backend:
	docker build . -f Dockerfile.backend -t dmitriko/djavue-backend:${TAG}

prepare-instance:
	venv/bin/ansible-playbook -i ./ops/inventory.ini ./ops/prepare-instance.yaml

deploy-backend:
	venv/bin/ansible-playbook -i ./ops/inventory.ini ./ops/deploy-backend.yaml

deploy-frontend:
	venv/bin/ansible-playbook -i ./ops/inventory.ini ./ops/deploy-frontend.yaml




