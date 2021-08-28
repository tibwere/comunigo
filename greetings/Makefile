PROTOFILE=encoding/proto/hello.proto

proto:
	protoc --go_out=. --go_opt=paths=source_relative \
    	--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		$(PROTOFILE)	

server: create-bin-directory
	go build -o ./bin/server ./cmd/server/main.go

client: create-bin-directory
	go build -o ./bin/client ./cmd/client/main.go

create-bin-directory:
	mkdir -p ./bin

clean:
	rm -rf ./bin

.PHONY: create-bin-directory clean

