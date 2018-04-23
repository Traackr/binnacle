# binnacle [![Release][release-image]][release-url] [![Build Status][travis-image]][travis-url]

`binnacle` is a command line tool used to interact with Kubernetes' Helm.

## Installation

A binary for various operating systems is available through [Github Releases][github-releases].  Download the appropriate archive, and extract into a directory within your PATH.

## Usage

For the full list of options:

```shell
binnacle --help
```

To see the version of `binnacle` you can use the following:

```shell
binnacle --version
```

## Getting Started

### Configuration File Format

## Development

To ease the entry of building `binnacle` there are two methods supported by the local Makefile.  The first is for a fully installed and configured [Go][go] (version 1.8+) environment on your machine, and the second requires only that docker be installed.

### Local Go Environment

You will first want to check out this repository into your GOPATH:

```script
mkdir -p "$GOPATH/src/github.com/traackr"
cd "$GOPATH/src/github.com/traackr"
git clone https://github.com/traackr/binnacle.git
cd binnacle
```

To compile a version of binnacle for your local machine you can run:

```script
make
```

This will generate a binary within the ./bin directory of the project.

To run the unit tests:

```script
make test-unit
```

To run the unit tests with coverage reports:

```script
make test-coverage
```

### Local Docker Environment

Using a local Docker environment for building runs the exact same commands as local development, they just happen to be run inside of the container.

To leverage the docker build environment you will first want to check out this repository into a directory of your choice.  In the example below there is an environment variable named `DEVELOPMENT` where all development files are stored.

```script
mkdir -p "$DEVELOPMENT/traackr"
cd "$DEVELOPMENT/traackr"
git clone https://github.com/traackr/binnacle.git
cd binnacle
```

To compile a version of binnacle for your local machine you can run:

```script
make docker-build
```

This will generate a binary within the ./bin directory of the project.

To run the unit tests:

```script
make docker-test-unit
```

To run the unit tests with coverage reports:

```script
make docker-test-coverage
```

[docker]: https://www.docker.com
[docker-compose]: https://docs.docker.com/compose/
[docker-golang]: https://hub.docker.com/_/golang/
[github-releases]: https://github.com/traackr/binnacle/releases
[go]: https://www.golang.org/
[release-url]: https://github.com/traackr/binnacle/releases/latest
[release-image]: https://img.shields.io/github/release/traackr/binnacle.svg
[travis-url]: https://travis-ci.org/traackr-binnacle
[travis-image]: https://travis-ci.org/traackr/binnacle.svg?branch=master
