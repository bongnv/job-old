.PHONY: install-tools lint test ci

test:
	@echo ">  Running tests..."
	go test -v -race ./...

lint:
	@echo "  Running go vet..."
	go vet ./...
	@echo "  Running golint.."
	golint -set_exit_status=1 ./...

install-tools:
	@echo ">  Installing tools..."
	go get -u golang.org/x/lint/golint

ci: install-tools lint test

