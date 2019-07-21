#############      builder       #############
FROM golang:1.12.5 AS builder

WORKDIR /go/src/github.com/leaderboard/Web-Cat-Leaderboard
COPY . .
RUN go get github.com/go-sql-driver/mysql
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go install \
  ./...

#############      apiserver     #############
FROM alpine:3.8 AS apiserver

RUN apk add --update bash curl tzdata

COPY --from=builder /go/bin/Web-Cat-Leaderboard /Web-Cat-Leaderboard

WORKDIR /

COPY form form/
ENTRYPOINT ["/Web-Cat-Leaderboard"]