FROM golang:alpine AS builder

ENV CGO_ENABLED 0

ENV GOOS linux

WORKDIR /build

ADD ./go.mod .

COPY . .

RUN go build -trimpath -o api ./api-gateway/api.go

COPY ./.env ./build/.env

FROM alpine

RUN apk update --no-cache && apk add --no-cache ca-certificates

WORKDIR /build

COPY --from=builder /build/api /build/api

EXPOSE 8080 8080

CMD ["/build/api"]