FROM alpine:3.19.1

RUN mkdir /app

COPY orderApp /app

RUN apk --no-cache add curl

CMD [ "/app/orderApp"]