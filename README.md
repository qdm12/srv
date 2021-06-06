# SRV

*srv is a small Go application to use in other projects as a base Docker image  serve static files over HTTP*

<img height="200" src="https://raw.githubusercontent.com/qdm12/srv/main/title.svg?sanitize=true">

[![Build status](https://github.com/qdm12/srv/workflows/CI/badge.svg)](https://github.com/qdm12/srv/actions?query=workflow%3ACI)
[![Docker Pulls](https://img.shields.io/docker/pulls/qmcgaw/srv.svg)](https://hub.docker.com/r/qmcgaw/srv)
[![Docker Stars](https://img.shields.io/docker/stars/qmcgaw/srv.svg)](https://hub.docker.com/r/qmcgaw/srv)
[![Image size](https://images.microbadger.com/badges/image/qmcgaw/srv.svg)](https://microbadger.com/images/qmcgaw/srv)
[![Image version](https://images.microbadger.com/badges/version/qmcgaw/srv.svg)](https://microbadger.com/images/qmcgaw/srv)

[![Join Slack channel](https://img.shields.io/badge/slack-@qdm12-yellow.svg?logo=slack)](https://join.slack.com/t/qdm12/shared_invite/enQtOTE0NjcxNTM1ODc5LTYyZmVlOTM3MGI4ZWU0YmJkMjUxNmQ4ODQ2OTAwYzMxMTlhY2Q1MWQyOWUyNjc2ODliNjFjMDUxNWNmNzk5MDk)
[![GitHub last commit](https://img.shields.io/github/last-commit/qdm12/srv.svg)](https://github.com/qdm12/srv/commits/main)
[![GitHub commit activity](https://img.shields.io/github/commit-activity/y/qdm12/srv.svg)](https://github.com/qdm12/srv/graphs/contributors)
[![GitHub issues](https://img.shields.io/github/issues/qdm12/srv.svg)](https://github.com/qdm12/srv/issues)

## Features

- Compatible with `amd64`, `386`, `arm64`, `arm32v7`, `arm32v6`, `ppc64le`, `s390x` and `riscv64` CPU architectures
- Runs without root as user ID `1000`
- Based on the `scratch` docker image for a tiny size of 9.11MB (uncompressed amd64 image)
- [Docker image tags and sizes](https://hub.docker.com/r/qmcgaw/srv/tags)
- Prometheus metrics available on `:9091`

## Setup

### As a container

We assume your static files are at `/yourpath` on your host.

1. Ensure the ownership of your static files matches the ones from the container:

    ```sh
    chown -R 1000 /yourpath
    ```

1. Run the container, bind mounting `/yourpath` to `/srv` as ready only:

    ```sh
    docker run -d -p 8000:8000/tcp -v /yourpath:/srv:ro qmcgaw/srv
    ```

    You can also use [docker-compose.yml](https://github.com/qdm12/srv/blob/main/docker-compose.yml) with:

    ```sh
    docker-compose up -d
    ```

1. You can update the image with `docker pull qmcgaw/srv` or use one of the [tags available](https://hub.docker.com/r/qmcgaw/srv/tags)

### As a base image

You can use it in your own project with:

```Dockerfile
FROM qmcgaw/srv
COPY --chown=1000 srv /srv
```

Or with a multi stage Dockerfile:

```Dockerfile
FROM someimage AS staticbuilder
COPY . .
RUN somecompiler --output-to=/tmp/build

FROM qmcgaw/srv
COPY --from=staticbuilder --chown=1000 /tmp/build /srv
```

### Environment variables

| Environment variable | Default | Possible values | Description |
| --- | --- | --- | --- |
| `HTTP_SERVER_ADDRESS` | `:8000` | Valid address | HTTP server listening address |
| `HTTP_SERVER_ROOT_URL` | `/` | URL path | HTTP server root URL |
| `HTTP_SERVER_LOG_REQUESTS` | `on` | `on` or `off` | Log requests and responses information |
| `HTTP_SERVER_SRV_FILEPATH` | `/srv` | Valid file path | File path to your static files directory |
| `HTTP_SERVER_ALLOWED_ORIGINS` | | CSV of addresses | Comma separated list of addresses to allow for CORS |
| `HTTP_SERVER_ALLOWED_HEADERS` | | CSV of HTTP header keys | Comma separated list of header keys to allow for CORS |
| `METRICS_SERVER_ADDRESS` | `:9090` | Valid address | Prometheus HTTP server listening address |
| `LOG_LEVEL` | `info` | `debug`, `info`, `warning`, `error` | Logging level |
| `HEALTH_SERVER_ADDRESS` | `127.0.0.1:9999` | Valid address | Health server listening address |
| `TZ` | `America/Montreal` | *string* | Timezone |

## Development

1. Setup your environment

    <details><summary>Using VSCode and Docker (easier)</summary><p>

    1. Install [Docker](https://docs.docker.com/install/)
       - On Windows, share a drive with Docker Desktop and have the project on that partition
       - On OSX, share your project directory with Docker Desktop
    1. With [Visual Studio Code](https://code.visualstudio.com/download), install the [remote containers extension](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers)
    1. In Visual Studio Code, press on `F1` and select `Remote-Containers: Open Folder in Container...`
    1. Your dev environment is ready to go!... and it's running in a container :+1: So you can discard it and update it easily!

    </p></details>

    <details><summary>Locally</summary><p>

    1. Install [Go](https://golang.org/dl/), [Docker](https://www.docker.com/products/docker-desktop) and [Git](https://git-scm.com/downloads)
    1. Install Go dependencies with

        ```sh
        go mod download
        ```

    1. Install [golangci-lint](https://github.com/golangci/golangci-lint#install)
    1. You might want to use an editor such as [Visual Studio Code](https://code.visualstudio.com/download) with the [Go extension](https://code.visualstudio.com/docs/languages/go). Working settings are already in [.vscode/settings.json](https://github.com/qdm12/srv/main/.vscode/settings.json).

    </p></details>

1. Commands available:

    ```sh
    # Build the binary
    go build cmd/app/main.go
    # Test the code
    go test ./...
    # Lint the code
    golangci-lint run
    # Build the Docker image
    docker build -t qmcgaw/srv .
    ```

1. See [Contributing](https://github.com/qdm12/srv/main/.github/CONTRIBUTING.md) for more information on how to contribute to this repository.

## License

This repository is under an [MIT license](https://github.com/qdm12/srv/main/license) unless otherwise indicated
