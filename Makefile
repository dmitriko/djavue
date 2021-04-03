# simple make file for django/vue  project

PIP = venv/bin/pip

venv:
	python3 -m venv venv
	$(PIP) install --upgrade -r requirements.txt

venv-update: venv
	$(PIP) install --upgrade -r requirements.txt

test: venv
	python manage.py test
