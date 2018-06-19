default: lint test

lint:
	golint ./...

test:
	go test -cover -v ./...

build:
	go build -o diceware main.go

win-build-all: win-build
	GOOS=windows GOARCH=amd64 go build -o diceware_64_bit.exe main.go

win-build:
	GOOS=windows GOARCH=386 go build -o diceware_32_bit.exe main.go

RM=rm

clean:
	$(RM) diceware
