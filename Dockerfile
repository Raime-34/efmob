# syntax=docker/dockerfile:1

FROM golang:1.25-alpine AS build

WORKDIR /src

RUN apk add --no-cache ca-certificates git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=0 GOOS=linux

RUN go build -trimpath -ldflags="-s -w" -o /out/subscriptionservice ./cmd/subscriptionservice

FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app

COPY --from=build /out/subscriptionservice /app/subscriptionservice
COPY --from=build /src/migrations /app/migrations

ENV DATABASE_URI=""

USER nonroot:nonroot

ENTRYPOINT ["/app/subscriptionservice"]
