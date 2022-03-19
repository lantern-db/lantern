# build stage
FROM golang:1.18.0-alpine3.15 AS builder
ADD . /src
RUN apk add git
RUN cd /src && go build -o /src/bin/viewer -v /src/viewer/cmd/

# final stage
FROM alpine:3.15

ADD ./viewer/static /app/viewer/static

WORKDIR /app
COPY --from=builder /src/bin/viewer /tmp/viewer
ENTRYPOINT /tmp/viewer