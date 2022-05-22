BINARY="search-engine"

build:
	go build -o ${BINARY}

build-linux:
	$ENV:CGO_ENABLED=0
	$ENV:GOOS="linux"
	go build -o ${BINARY}

.PHONY: build build-linux