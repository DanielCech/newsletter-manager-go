ARG BUILDER_IMAGE_VERSION=1.20
ARG RUNTIME_IMAGE_VERSION=3.16

FROM golang:${BUILDER_IMAGE_VERSION} as builder
ARG APP_VERSION=0.0.0
WORKDIR /go/src/go.strv.io/newsletter-manager-go
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-X main.version=${APP_VERSION}" -o bin/api cmd/api/main.go

FROM alpine:${RUNTIME_IMAGE_VERSION} as runtime
COPY --from=builder /go/src/go.strv.io/newsletter-manager-go/bin/api /usr/local/bin/api
ENTRYPOINT ["/usr/local/bin/api"]
