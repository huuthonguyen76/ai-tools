generate-docs:
	cd be && swag init --dir cmd/ -o docs

run-be:
	cd be && go run cmd/*.go