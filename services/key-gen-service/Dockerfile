ARG GO_VERSION=1.23.0
FROM golang:${GO_VERSION} AS build

WORKDIR /src
RUN --mount=type=cache,target=/go/pkg/mod/ \
  --mount=type=bind,source=go.mod,target=go.mod \
  --mount=type=bind,source=go.sum,target=go.sum \
  go mod download -x

ARG VERSION=
ARG COMMIT=
ARG DATE=
RUN --mount=type=cache,target=/go/pkg/mod \
  --mount=type=bind,target=. \
  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \ 
  go build -ldflags="-s -w -X 'main.version=${VERSION}' -X 'main.commit=${COMMIT}' -X 'main.date=${DATE}'" \
  -o /bin/key-gen .

FROM alpine:3.20.0 AS release

RUN --mount=type=cache,target=/var/cache/apk/ \
  apk --update add \
  ca-certificates \
  tzdata \
  && \
  update-ca-certificates

ARG UID=10001
RUN adduser -H -D \
  --uid "${UID}" appuser
USER appuser

COPY --from=build /bin/key-gen /usr/bin/

EXPOSE 8080 9090

ENTRYPOINT [ "/usr/bin/key-gen" ]
