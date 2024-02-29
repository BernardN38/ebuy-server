FROM alpine:latest

RUN mkdir /app

COPY target/x86_64-unknown-linux-musl/release/user-service /app

RUN apk --no-cache add curl
# RUN apk --no-cache add go
CMD [ "/app/user-service"]