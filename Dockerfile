FROM golang:1.23-alpine AS builder

WORKDIR /usr/local/src

RUN addgroup -S 1001 && adduser -S crs -G 1001

RUN apk --no-cache add bash git make gcc gettext musl-dev

ADD ["go.mod", "go.sum", "./"]

RUN go mod download

COPY . .

RUN go build -o ./bin/server cmd/server/main.go

FROM alpine AS runner

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /usr/local/src/bin/server /

USER crs

EXPOSE 3000

CMD ["/server"]