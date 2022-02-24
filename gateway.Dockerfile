# build stage
FROM golang:1.17.7-alpine3.15 AS builder
ADD . /src
RUN cd /src && go build -o /src/bin/lantern-gateway -v /src/gateway/cmd/

# final stage
FROM alpine:3.15
ENV LANTERN_HOST=localhost
ENV LANTERN_PORT=6380
ENV GATEWAY_PORT=8080

WORKDIR /app
COPY --from=builder /src/bin/lantern-gateway /tmp/lantern-gateway
ENTRYPOINT /tmp/lantern-gateway