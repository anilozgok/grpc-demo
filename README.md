# gRPC Demo 

To genereate the gRPC code, run the following command:

    protoc --go_out=. --go-grpc_out=. proto/demo.proto

To run the server, run the following command:

    go run server/main.go

To run the client, run the following command:

    go run client/main.go