run:
	go run main.go

proto:
	protoc --proto_path=pb --go_out=pb --go_opt=paths=source_relative \
	--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	pb/*.proto

clean:
	rm pb/*.pb.go

# grpcurl execution
create:
	grpcurl --plaintext -d '{"user": "customeruser", "pass": "secret", "detail": {"name": "SMK Walisongo", "email": "walisongo@gmail.com", "wa": "083879154310"}}' localhost:5012 CustomerService.Create

findOne:
	grpcurl --plaintext -d '{"user": "customeruser", "pass": "secret"}' localhost:5012 CustomerService.FindOne

find:
	grpcurl --plaintext localhost:5012 CustomerService.FindAll
	
findSearch:
	grpcurl --plaintext -d '{"search": "1671196175094202250", "status": "deleted"}' localhost:5012 CustomerService.FindAll
	