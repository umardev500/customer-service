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

changeStatus:
	grpcurl --plaintext -d '{"customer_id": "1671196175094202250", "status": "accepted"}' localhost:5012 CustomerService.ChangeStatus

updateDetail:
	grpcurl --plaintext -d '{"customer_id": "1671196175094202250", "detail": {"npsn": "1923823", "name": "SMK Walisongo", "email": "walisongo@gmail.com", "wa": "083879154310", "type": "private", "level": "High School", "about": "This is about the institution", "location": {"address": "Jl. Labuan km.12", "village": "menes", "district": "menes", "city": "pandeglang", "province": "banten", "postal_code": "42262"}}}' localhost:5012 CustomerService.UpdateDetail
	
findSearch:
	grpcurl --plaintext -d '{"search": "1671196175094202250", "status": "deleted"}' localhost:5012 CustomerService.FindAll

delete:
	grpcurl --plaintext -d '{"customer_id": "1671079334676443995", "hard": true}' localhost:5012 CustomerService.Delete
	