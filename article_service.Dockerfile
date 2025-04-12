FROM golang:alpine AS builder

ENV CGO_ENABLED 0

ENV GOOS linux

WORKDIR /build

ADD ./go.mod .

COPY . .

RUN go build -trimpath -o article ./article-service/article.go

FROM alpine

RUN apk update --no-cache && apk add --no-cache ca-certificates

WORKDIR /build

COPY --from=builder /build/article /build/article

EXPOSE 50052 50052

CMD ["/build/article"]