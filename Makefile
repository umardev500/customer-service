run:
	go run main.go

proto:
	protoc --go_out=. --go_opt=paths=source_relative \
	--go-grpc_out=. --go-grpc_opt=paths=source_relative \
	pb/*.proto

clean:
	rm pb/*.pb.go

# grpcurl execution
create:
	grpcurl --plaintext -d '{"user": "customeruser", "pass": "secret", "detail": {"name": "SMK Walisongo", "email": "walisongo@gmail.com", "wa": "083879154310"}}' localhost:5012 CustomerService.Create

find:
	grpcurl --plaintext -d '{"customer_id": "1667292823233"}' localhost:5012 CustomerService.Find

findAll:
	grpcurl --plaintext localhost:5012 CustomerService.FindAll

changeStatus:
	grpcurl --plaintext -d '{"customer_id": "1671196175094202250", "status": "accepted"}' localhost:5012 CustomerService.ChangeStatus

updateDetail:
	grpcurl --plaintext -d '{"customer_id": "1667292823233", "detail": {"npsn": "150", "name": "SMK Walisongo", "email": "walisongo@gmail.com", "wa": "+62 83879154310", "type": "private", "level": "High School", "about": "This is about the institution", "logo": "/uploads/images/avatars/avatar.png", "location": {"address": "Jl. Labuan km.12 Ciputri", "village": "Menes", "district": "Menes", "city": "pandeglang", "province": "banten", "postal_code": "42262"}}}' localhost:5012 CustomerService.UpdateDetail

updateDetailShort:
	grpcurl --plaintext -d '{"customer_id": "1667292823233", "detail": {"name": "150"}}' localhost:5012 CustomerService.UpdateDetail

updateCreds:
	grpcurl --plaintext -d '{"user": "walisongo", "pass": "walisongopass", "new_pass": "walisongo155pass"}' localhost:5012 CustomerService.UpdateCreds
	
findSearch:
	grpcurl --plaintext -d '{"search": "", "status": "deleted"}' localhost:5012 CustomerService.FindAll

delete:
	grpcurl --plaintext -d '{"customer_id": "1671079334676443995", "hard": true}' localhost:5012 CustomerService.Delete

setExp:
	grpcurl --plaintext -d '{"customer_id": "16672928232332", "exp_time": 1}' localhost:5012 CustomerService.SetExp

login:
	grpcurl --plaintext -d '{"user": "walisongo", "pass": "walisongopass"}' localhost:5012 CustomerService.Login
	