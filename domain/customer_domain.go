package domain

import (
	"context"
	"customer/pb"
)

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

type CustomerLoginResponse struct {
	CustomerId string `bson:"customer_id"`
	User       string `bson:"user"`
}

type CustomerUsecase interface {
	Save(ctx context.Context, req *pb.CustomerCreateRequest) error
	FindOne(ctx context.Context, req *pb.CustomerFindOneRequest) (customer *pb.Customer, err error)
	FindAll(ctx context.Context, req *pb.CustomerFindAllRequest) (customers *pb.CustomerFindAllResponse, err error)
	ChangeStatus(ctx context.Context, req *pb.CustomerChangeStatusRequest) (affected bool, err error)
	UpdateDetail(ctx context.Context, req *pb.CustomerUpdateDetailRequest) (affected bool, err error)
	Delete(ctx context.Context, req *pb.CustomerDeleteRequest) (affected bool, err error)
	SetExp(ctx context.Context, req *pb.CustomerSetExpRequest) (affected bool, err error)
	Login(ctx context.Context, req *pb.CustomerLoginRequest) (res *pb.CustomerLoginResponse, err error)
}

type CustomerRepository interface {
	Save(ctx context.Context, req *pb.CustomerCreateRequest, generatedId string, createdTime int64) error
	FindOne(ctx context.Context, req *pb.CustomerFindOneRequest) (customer *pb.Customer, err error)
	FindAll(ctx context.Context, req *pb.CustomerFindAllRequest) (customers *pb.CustomerFindAllResponse, err error)
	ChangeStatus(ctx context.Context, req *pb.CustomerChangeStatusRequest, updatedTime int64) (affected bool, err error)
	UpdateDetail(ctx context.Context, req *pb.CustomerUpdateDetailRequest, updatedTime int64) (affected bool, err error)
	Delete(ctx context.Context, req *pb.CustomerDeleteRequest, deletedTime int64) (affected bool, err error)
	SetExp(ctx context.Context, req *pb.CustomerSetExpRequest, updatedTime int64) (affected bool, err error)
	Login(ctx context.Context, req *pb.CustomerLoginRequest) (res *pb.CustomerLoginResponse, err error)
}
