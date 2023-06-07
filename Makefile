all: build

build:
	# GOOS=linux GOARCH=amd64 go build  -o bin/nb-extractor ./cmd/nb-extractor
	go build  -o bin/nb-extractor ./cmd/nb-extractor