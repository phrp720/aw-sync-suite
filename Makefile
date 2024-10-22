run:
	@go run main.go

dev-run:
	@go run main.go -asService=false -awUrl=http://localhost:5600 -prometheusUrl=http://localhost:9090 -userID=DevUser

build:
	@go build -o bin/agent main.go
	@cp -r .env bin/.env

test:
	@go test -v ./...

clean:
	@rm -rf bin

format:
	@gofmt -s -w .