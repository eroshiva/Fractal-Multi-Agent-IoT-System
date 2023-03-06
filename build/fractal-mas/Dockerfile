# this is a minimal image, which needs git (and curl) to be installed
FROM golang:1.20.1-alpine

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates git curl

WORKDIR /app

ADD bin /bin/

CMD ["/bin/sh"]