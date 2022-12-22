FROM alpine:latest

COPY ./build/inHabrBot /app/
COPY configs /app/configs

WORKDIR /app

ENTRYPOINT ["./inHabrBot"]