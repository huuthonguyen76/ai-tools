test:
	poetry run pytest tests/unit_test/ --capture=no

test-integration:
	poetry run pytest tests/integration/ --capture=no
