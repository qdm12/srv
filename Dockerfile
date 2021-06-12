# Sets linux/amd64 in case it's not injected by older Docker versions
ARG BUILDPLATFORM=linux/amd64

ARG ALPINE_VERSION=3.13
ARG GO_VERSION=1.16

FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS base
RUN apk --update add git g++
ENV CGO_ENABLED=0
ARG GOLANGCI_LINT_VERSION=v1.40.1
RUN go get github.com/golangci/golangci-lint/cmd/golangci-lint@${GOLANGCI_LINT_VERSION}
COPY --from=qmcgaw/xcputranslate:v0.4.0 /xcputranslate /usr/local/bin/xcputranslate
WORKDIR /tmp/gobuild
COPY go.mod go.sum ./
RUN go mod download
COPY cmd/ ./cmd/
COPY internal/ ./internal/

FROM base AS test
# Note on the go race detector:
# - we set CGO_ENABLED=1 to have it enabled
# - we installed g++ in the base stage to support the race detector
ENV CGO_ENABLED=1

FROM base AS lint
COPY .golangci.yml ./
RUN golangci-lint run --timeout=10m

FROM base AS tidy
RUN git init && \
    git config user.email ci@localhost && \
    git config user.name ci && \
    git add -A && git commit -m ci && \
    sed -i '/\/\/ indirect/d' go.mod && \
    go mod tidy && \
    git diff --exit-code -- go.mod

FROM base AS build
ARG TARGETPLATFORM
ARG VERSION=unknown
ARG BUILD_DATE="an unknown date"
ARG COMMIT=unknown
RUN GOARCH="$(xcputranslate -targetplatform=${TARGETPLATFORM} -field arch)" \
    GOARM="$(xcputranslate -targetplatform=${TARGETPLATFORM} -field arm)" \
    go build -trimpath -ldflags="-s -w \
    -X 'main.version=$VERSION' \
    -X 'main.buildDate=$BUILD_DATE' \
    -X 'main.commit=$COMMIT' \
    " -o app cmd/app/main.go

FROM --platform=$BUILDPLATFORM alpine:${ALPINE_VERSION} AS alpine

FROM scratch
USER 1000
ENTRYPOINT ["/app"]
EXPOSE 8000/tcp
HEALTHCHECK --interval=10s --timeout=5s --start-period=5s --retries=2 CMD ["/app","healthcheck"]
ENV HTTP_SERVER_ADDRESS=:8000 \
    HTTP_SERVER_ROOT_URL=/ \
    HTTP_SERVER_LOG_REQUESTS=on \
    FILEPATH_SRV=./srv \
    FILEPATH_WORK=/tmp/srv \
    METRICS_SERVER_ADDRESS=:9090 \
    LOG_LEVEL=info \
    HEALTH_SERVER_ADDRESS=127.0.0.1:9999 \
    TZ=America/Montreal
COPY --from=alpine --chown=1000 /srv /srv
COPY --from=alpine --chown=1000 /srv /tmp/srv
ARG VERSION=unknown
ARG BUILD_DATE="an unknown date"
ARG COMMIT=unknown
LABEL \
    org.opencontainers.image.authors="quentin.mcgaw@gmail.com" \
    org.opencontainers.image.version=$VERSION \
    org.opencontainers.image.created=$BUILD_DATE \
    org.opencontainers.image.revision=$COMMIT \
    org.opencontainers.image.url="https://github.com/qdm12/srv" \
    org.opencontainers.image.documentation="https://github.com/qdm12/srv/blob/main/README.md" \
    org.opencontainers.image.source="https://github.com/qdm12/srv" \
    org.opencontainers.image.title="srv" \
    org.opencontainers.image.description="srv is a small Go application to use as a container or as a base Docker image in other projects to serve static files over HTTP"
VOLUME [ "/tmp/srv" ]
COPY --from=build --chown=1000 /tmp/gobuild/app /app