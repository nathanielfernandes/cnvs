gen_proto:
	protoc --go_out=./canvas ./protobufs/canvas.proto

run:
	go run ./api

build:
	go build -o ./bin/canvas ./api