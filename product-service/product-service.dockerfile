FROM alpine:latest

RUN mkdir /app

COPY productApp /app

RUN apk --no-cache add curl

CMD [ "/app/productApp"]