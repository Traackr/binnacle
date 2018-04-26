# binnacle [![Release][release-image]][release-url] [![Build Status][travis-image]][travis-url]

`binnacle` is an opinionated automation tool used to interact with Kubernetes' [Helm][helm].  By offering a single file to manage one or many charts, you can easily control all aspects of your Helm managed applications.

`binnacle` is similar in nature to [Helmfile][helmfile] with a slightly different appraoch to managing Helm Charts.

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

Configuration files can be written in YAML, TOML or JSON.

```yaml
---
# charts takes a list of chart configurations
charts:
    # This is the name of the chart
  - name: concourse
    # This is the namespace into which the chart is launched
    namespace: apps
    # This is the name for the release of this chart
    release: apps-concourse
    # This is the name of the repository within which the helm chart is located
    repo: stable
    # This determines if the release is created or removed. Default: present Options: absent, present
    state: present
    # Any data under values are passed to Helm to configure the given chart
    values:
      image: concourse/concourse
      imageTag: "3.10.0"
    # This is the version of the Helm chart.  If this is omitted, the latest is used.
    version: 1.3.1

# repositories takes a list of repository configurations
repositories:
    # This is the name of the repository
  - name: stable
    # This is the URL of the repository
    url: https://kubernetes-charts.storage.googleapis.com
    # This determines if the repository is created or removed. Default: present Options: absent, present
    state: present
```

### Commands

Documentation for all of the commands within `binnacle` are available [here][commands].

## Development

To ease the entry of building `binnacle` there are two methods supported by the local Makefile.  The first is for a fully installed and configured [Go][go] (version 1.8+) environment on your machine, and the second requires only that docker be installed.

### Local Go Environment

You will first want to check out this repository into your GOPATH:

```script
mkdir -p "$GOPATH/src/github.com/Traackr"
cd "$GOPATH/src/github.com/Traackr"
git clone https://github.com/Traackr/binnacle.git
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

Using a local [Docker][docker] environment for building runs the exact same commands as local development, they just happen to be run inside of the container.

To leverage the docker build environment you will first want to check out this repository into a directory of your choice.  In the example below there is an environment variable named `DEVELOPMENT` where all development files are stored.

```script
mkdir -p "$DEVELOPMENT/Traackr"
cd "$DEVELOPMENT/Traackr"
git clone https://github.com/Traackr/binnacle.git
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

[commands]: docs/commands/binnacle.md
[docker]: https://www.docker.com
[github-releases]: https://github.com/Traackr/binnacle/releases
[go]: https://www.golang.org/
[helm]: https://helm.sh/
[helmfile]: https://github.com/roboll/helmfile
[release-url]: https://github.com/Traackr/binnacle/releases/latest
[release-image]: https://img.shields.io/github/release/Traackr/binnacle.svg
[travis-url]: https://travis-ci.org/Traackr/binnacle
[travis-image]: https://travis-ci.org/Traackr/binnacle.svg?branch=master
