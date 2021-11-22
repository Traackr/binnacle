.PHONY: build clean test

build:
	@scripts/build.sh

clean:
	@rm -rf bin pkg

test:
	@go test -v ./...