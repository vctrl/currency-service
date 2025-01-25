package main

//go:generate protoc --proto_path=proto --go_out=$(PROTO_OUT) --go-grpc_out=$(PROTO_OUT) $(PROTO_SRC)
