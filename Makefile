build:
		@go build -o bin/gobanking

run: build 
		@./bin/gobanking

test:
		@go test -v ./..