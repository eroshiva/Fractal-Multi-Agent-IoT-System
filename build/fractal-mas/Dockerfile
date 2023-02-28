# FIXME: fix this Dockerfile
# this is a minimal destribution of Debian - consider to use Go image for the tests...
#FROM debian:bullseye-slim
# this is a slim image, which needs git (and curl) to be installed
FROM golang:1.20.1-alpine

WORKDIR /

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates git curl

ADD bin /bin/

CMD ["/bin/sh"]