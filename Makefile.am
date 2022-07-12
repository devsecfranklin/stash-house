CLOCAL_AMFLAGS = -I config/m4 -I aclocal
ACLOCAL_AMFLAGS = -I config/m4 -I aclocal

python:
	@$(PYTHON) -m venv _build
	( \
		source _build/bin/activate; \
		_build/bin/python -m pip install --no-index --trusted-host storage1.lab.bitsmasher.net --find-links=http://storage1.lab.bitsmasher.net/pypi tox; \
		_build/bin/python -m pip install --trusted-host storage1.lab.bitsmasher.net --find-links=http://storage1.lab.bitsmasher.net/pypi --upgrade pip; \
		_build/bin/python -m pip install --trusted-host storage1.lab.bitsmasher.net --find-links=http://storage1.lab.bitsmasher.net/pypi -r src/requirements.txt; \
	)
