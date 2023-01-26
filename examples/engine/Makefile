.PHONY: compile

compile:
	protoc -Iproto/ --go_out=. --go_opt=module=github.com/SergeyParamoshkin/rebrainme/engine --go-grpc_out=. --go-grpc_opt=module=github.com/SergeyParamoshkin/rebrainme/engine proto/*.proto

run:
	go run cmd/main.go