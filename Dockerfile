# build stage
FROM golang:1.16.3-alpine3.13 AS builder
ADD . /src
RUN cd /src && go build -o /src/bin/lantern-server -v /src/server/cmd/server.go

# final stage
FROM alpine
ENV LANTERN_FLUSH_INTERVAL=60
ENV LANTERN_PORT=6380
ENV LANTERN_TTL=180

WORKDIR /app
COPY --from=builder /src/bin/lantern-server /tmp/lantern-server
ENTRYPOINT /tmp/lantern-server