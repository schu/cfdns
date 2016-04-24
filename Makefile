all: build

build:
	go build -o bin/cfdns \
		main.go
