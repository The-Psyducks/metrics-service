# Build stage
FROM golang:1.23.0-alpine3.20 AS builder

WORKDIR /home/app

COPY /server/go.mod /server/go.sum ./
RUN go mod download

COPY /server ./

RUN go build -o twitsnap ./main.go


FROM alpine:3.20

WORKDIR /home/app

COPY --from=builder /home/app/twitsnap ./

ENTRYPOINT ["./twitsnap"]