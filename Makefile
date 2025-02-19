run:
	go run .

test:
	go test -v ./...

mockgen:
	go generate ./...

.PHONY: run test mockgen