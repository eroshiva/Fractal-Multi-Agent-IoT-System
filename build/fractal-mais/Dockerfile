# This is inspired by:
# https://github.com/onosproject/onos-e2t/blob/master/build/onos-e2t/Dockerfile

# this is a minimal image, which needs git (and curl) to be installed
FROM golang:1.20.4-alpine as builder

# installing dependencies
RUN apk add --update && apk add --update --no-cache make

# building a binary
ENV GO111MODULE=on
ARG MAKE_TARGET=build

COPY Makefile go.mod go.sum /fmais/
COPY cmd/ /fmais/cmd/
COPY pkg/ /fmais/pkg/
COPY vendor/ /fmais/vendor/

WORKDIR /fmais

RUN make ${MAKE_TARGET}
RUN mkdir -p /fmais/data && mkdir -p /fmais/figures

# building a small image
FROM scratch

# Copy our static executable
COPY --from=builder /fmais/build/_output/fractal-mais /usr/local/bin/fractal-mais
COPY --from=builder /fmais/data /usr/local/bin/data
COPY --from=builder /fmais/figures /usr/local/bin/figures
WORKDIR /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/fractal-mais"]
