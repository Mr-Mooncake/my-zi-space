.PHONY: run build test clean

run:
	go run .

build:
	go build -o bin/my-zi-space .

test:
	go test ./...

clean:
	rm -rf bin
