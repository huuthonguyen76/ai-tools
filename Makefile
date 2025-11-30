generate-docs:
	cd be && swag init --dir cmd/ -o docs

run-be:
	cd be && go run cmd/*.go

test-be:
	cd be && go test ./... -v

test-be-coverage:
	cd be && go test ./... -cover

test-be-specific:
	cd be/tests && go test -v

install-ui:
	cd ui && pip install -r requirements.txt

run-ui:
	cd ui && streamlit run app.py

run-all:
	make -j 2 run-be run-ui