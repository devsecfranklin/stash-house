# --- Python Virtual Environment Setup ---
VENV = .venv
VENV_BIN = $(VENV)/bin
PYTHON = $(VENV_BIN)/python3
PIP = $(VENV_BIN)/pip

# Rule to create the venv and install dependencies from both files
$(VENV)/bin/activate: requirements.txt test/requirements-test.txt
	@echo "==> Hydrating Environment with Core and Test Dependencies..."
	test -d $(VENV) || python3 -m venv $(VENV)
	$(PIP) install --upgrade pip
	$(PIP) install -r requirements.txt
	@if [ -f "test/requirements-test.txt" ]; then \
		echo "==> Installing Test Dependencies..."; \
		$(PIP) install -r test/requirements-test.txt; \
	fi
	@touch $(VENV)/bin/activate

python: $(VENV)/bin/activate
	@echo "==> Environment ready. Source with: source $(VENV)/bin/activate"

test: python
	@echo "==> Initiating Stargate Distributed Test Suite..."
	@# Setting PYTHONPATH ensures 'internal' is found inside 'src/'
	PYTHONPATH=src $(PYTHON) -m pytest test/

.PHONY: python test
