# build stage
FROM golang:1.18.0-alpine3.15 AS builder
ADD . /src
RUN apk add git
RUN cd /src && go build -o /src/bin/lantern-server -v /src/server/cmd/

# final stage
FROM alpine:3.15
ENV LANTERN_FLUSH_INTERVAL=60
ENV LANTERN_PORT=6380

WORKDIR /app
COPY --from=builder /src/bin/lantern-server /tmp/lantern-server
ENTRYPOINT /tmp/lantern-server