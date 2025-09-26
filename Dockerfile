FROM --platform=linux/amd64 debian:stable-slim

RUN apt-get update && apt-get install -y ca-certificates

ADD weather-api /usr/bin/weather-api

COPY .env .

CMD ["weather-api"]