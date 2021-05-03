# build stage
FROM golang:1.16.3-alpine3.13 AS builder
ADD . /src
RUN cd /src && go build -o /src/bin/server -v /src/cmd/server.go

# final stage
FROM alpine
ENV LANTERNE_FLUSH_INTERVAL=60
ENV LANTERNE_PORT=6380
ENV LANTERNE_TTL=180

WORKDIR /app
COPY --from=builder /src/bin/server /tmp/server
ENTRYPOINT /tmp/server