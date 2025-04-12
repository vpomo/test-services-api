FROM golang:alpine AS builder

ENV CGO_ENABLED 0

ENV GOOS linux

WORKDIR /build

ADD ./go.mod .

COPY . .

RUN go build -trimpath -o user ./user-service/user.go

COPY ./.env ./build/.env

FROM alpine

RUN apk update --no-cache && apk add --no-cache ca-certificates

WORKDIR /build

COPY --from=builder /build/user /build/user

COPY --from=builder ./build/.env /build/.env

EXPOSE 50051 50051

CMD ["/build/user"]