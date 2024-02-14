FROM alpine:latest

RUN mkdir /app

COPY orderApp /app

RUN apk --no-cache add curl

CMD [ "/app/orderApp"]