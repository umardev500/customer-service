package domain

import "customer/pb"

type CustomerLocation struct {
	Address    string `bson:"address"`
	Village    string `bson:"village"`
	District   string `bson:"district"`
	City       string `bson:"city"`
	Province   string `bson:"province"`
	PostalCode string `bson:"postal_code"`
}

type CustomerDetail struct {
	Npsn     string            `bson:"npsn"`
	Name     string            `bson:"name"`
	Email    string            `bson:"email"`
	Wa       string            `bson:"wa"`
	Type     string            `bson:"type"`
	Level    string            `bson:"level"`
	About    string            `bson:"about"`
	Location *CustomerLocation `bson:"location"`
}

type Customer struct {
	CustomerId string          `bson:"customer_id"`
	User       string          `bson:"user"`
	Pass       string          `bson:"pass"`
	Detail     *CustomerDetail `bson:"detail"`
	Status     string          `bson:"status"`
	ExpUntil   int64           `bson:"exp_until"`
	CreatedAt  int64           `bson:"created_at"`
	UpdatedAt  int64           `bson:"updated_at"`
	DeletedAt  int64           `bson:"deleted_at"`
}

type CustomerUsecase interface {
	Save(req *pb.CustomerCreateRequest) error
	FindOne(req *pb.CustomerFindOneRequest) (customer *pb.Customer, err error)
}

type CustomerRepository interface {
	Save(req *pb.CustomerCreateRequest, generatedId string, createdTime int64) error
	FindOne(req *pb.CustomerFindOneRequest) (customer *pb.Customer, err error)
}
