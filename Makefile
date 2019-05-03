.PHONY: test
test: $(SRC)
	go test -v ./...

lint: $(SRC)
	golangci-lint run -v