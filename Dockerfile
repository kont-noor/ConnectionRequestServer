FROM golang:1.24.0-alpine3.21 AS builder

WORKDIR /app

RUN addgroup -S 1001 && adduser -S crs -G 1001

RUN apk --no-cache add bash git make gcc gettext musl-dev

ADD ["go.mod", "go.sum", "./"]

RUN --mount=type=cache,target=/go-cache \
    --mount=type=cache,target=/gomod-cache \
    go mod download

COPY . .

RUN --mount=type=cache,target=/go-cache \
    --mount=type=cache,target=/gomod-cache \
    go build \
        -ldflags="-linkmode external -extldflags -static" \
        -o ./bin/server cmd/server/main.go

FROM scratch AS runner

WORKDIR /app

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /app/bin/server /app/server

USER crs

EXPOSE 3000

CMD ["/app/server"]