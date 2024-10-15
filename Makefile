run:
	@go run main.go

build:
	@go build -o bin/agent main.go
	@cp -r .env bin/.env

test:
	@go test -v ./...

clean:
	@rm -rf bin