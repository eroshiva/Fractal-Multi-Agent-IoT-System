# this is a minimal image, which needs git (and curl) to be installed
FROM golang:1.20.2-alpine

RUN apk add --update && apk add --update --no-cache ca-certificates git curl
WORKDIR /app
COPY ./build/_output/fractal-mas /app/fractal-mas

#CMD ["/bin/sh"] /app/fractal-mas --example