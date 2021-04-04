# simple make file for django/vue  project

PIP = venv/bin/pip

venv:
	python3 -m venv venv

venv-update: venv
	$(PIP) install --upgrade -r requirements.txt
	$(PIP) install --upgrade -r requirements-dev.txt

test: venv
	python manage.py test

build-backend:
	docker build . -f Dockerfile.backend -t dmitriko/djavue-backend:latest
