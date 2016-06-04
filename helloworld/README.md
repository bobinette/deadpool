# Hello, World!

Protoc command:
```bash
protoc -I ./protos ./protos/*.proto --go_out=Mgoogle/protobuf/timestamp.proto=github.com/golang/protobuf/ptypes/timestamp,plugins=grpc:./protos
```
