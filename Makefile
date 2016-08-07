default: lint test

lint:
	golint ./...

test:
	go test -cover -v ./...

build:
	go build -o diceware main.go

RM=rm

clean:
	$(RM) diceware
