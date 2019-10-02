FROM alpine:3.7

MAINTAINER hakkim.badhusha@wipro.com

WORKDIR /

WORKDIR /app

COPY . .

WORKDIR /app/cmd


ENTRYPOINT ["/app/user-service"]