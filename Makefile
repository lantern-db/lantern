
./graph/model/mock/vertex.go: ./graph/model/vertex.go
	mockgen -source=./graph/model/vertex.go -destination=./graph/model/mock/vertex.go

./graph/model/mock/edge.go: ./graph/model/edge.go
	mockgen -source=./graph/model/edge.go -destination=./graph/model/mock/edge.go

./server/cmd/wire_gen.go: ./server/cmd/wire.go
	wire ./server/cmd/

./gateway/cmd/wire_gen.go: ./gateway/cmd/wire.go
	wire ./gateway/cmd/

.PHONY: test
test: ./graph/model/mock/vertex.go ./graph/model/mock/edge.go ./server/cmd/wire_gen.go ./gateway/cmd/wire_gen.go
	go mod tidy
	go build ./...
	go test ./...

.PHONY: build
build: server.Dockerfile test
	docker build -t lantern-server -f server.Dockerfile .
	docker tag lantern-server piroyoung/lantern-server:local
	docker tag lantern-server piroyoung/lantern-server:latest-alpha
	docker build -t lantern-gateway -f gateway.Dockerfile .
	docker tag lantern-gateway piroyoung/lantern-gateway:local
	docker tag lantern-gateway piroyoung/lantern-gateway:latest-alpha
	docker build -t lantern-viewer -f viewer.Dockerfile .
	docker tag lantern-viewer piroyoung/lantern-viewer:local
	docker tag lantern-viewer piroyoung/lantern-viewer:latest-alpha

.PHONY: push
push: build
	docker push piroyoung/lantern-server:latest-alpha
	docker push piroyoung/lantern-gateway:latest-alpha
	docker push piroyoung/lantern-viewer:latest-alpha

.PHONY: run
run: build
	docker-compose up


clean:
	rm ./graph/model/mock/*
.PHONY: clean