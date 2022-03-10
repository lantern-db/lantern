# build stage
FROM golang:1.17.7-alpine3.15 AS builder
ADD . /src
RUN cd /src && go build -o /src/bin/viewer -v /src/viewer/cmd/

# final stage
FROM alpine:3.15

ADD ./viewer/static /app/viewer/static

WORKDIR /app
COPY --from=builder /src/bin/viewer /tmp/viewer
ENTRYPOINT /tmp/viewer