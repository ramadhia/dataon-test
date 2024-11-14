# ==========================================
# 1st Stage
# ==========================================
FROM golang:1.21 AS builder
LABEL maintainer="rafli.ramadhia@gmail.com"

## Set the working directory
WORKDIR /app

## Copy source
COPY . .

## Compile
RUN make build

# ==========================================
# 2nd Stage
# ==========================================
FROM alpine:latest

ENV APP_NAME=be

WORKDIR /app

## Add ssl cert
RUN apk add --update --no-cache ca-certificates

## Add timezone data
RUN apk --no-cache add tzdata

## Copy binary file from 1st stage
COPY --from=builder /app/bin/* ./

## Copy migration files
#COPY ./database ./database

CMD ["./cli", "server"]
