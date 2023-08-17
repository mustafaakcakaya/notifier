FROM golang:1.20-alpine AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY *.go ./

RUN go build -o /notifier

FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build ./ /notifier

USER nonroot:nonroot

ENTRYPOINT ["/notifier"]