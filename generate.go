package main

//go:generate bash -c "mkdir -p pkg/pb"
//go:generate bash -c "protoc --go_out=pb --go_opt=Mproto/reverse.proto=github.com/SergeyParamoshkin/rebrainme --go-grpc_out=pb --go-grpc_opt=Mproto/reverse.proto=github.com/SergeyParamoshkin/rebrainme --proto_path=proto  proto/*.proto"
