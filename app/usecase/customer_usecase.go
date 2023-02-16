package usecase

import (
	"context"
	"customer/domain"
	"customer/pb"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

// CustomerUsecase defines the use case for managing Customers.
type CustomerUsecase struct {
	// repository is the underlying repository for storing Customers.
	repository domain.CustomerRepository
}

// NewCustomerUsecase creates a new CustomerUsecase with the given repository.
func NewCustomerUsecase(repo domain.CustomerRepository) domain.CustomerUsecase {
	return &CustomerUsecase{
		repository: repo,
	}
}

// Template
// func (pu *CustomerUsecase) {}

func (c *CustomerUsecase) Login(ctx context.Context, req *pb.CustomerLoginRequest) (res *pb.CustomerLoginResponse, err error) {
	res, err = c.repository.Login(ctx, req)
	return
}

func (c *CustomerUsecase) SetExp(ctx context.Context, req *pb.CustomerSetExpRequest) (affected bool, err error) {
	updatedTime := time.Now().UTC().Unix()
	affected, err = c.repository.SetExp(ctx, req, updatedTime)

	return
}

func (c *CustomerUsecase) Delete(ctx context.Context, req *pb.CustomerDeleteRequest) (affected bool, err error) {
	deletedTime := time.Now().UTC().Unix()
	affected, err = c.repository.Delete(ctx, req, deletedTime)

	return
}

func (c *CustomerUsecase) UpdateDetail(ctx context.Context, req *pb.CustomerUpdateDetailRequest) (affected bool, err error) {
	updatedTime := time.Now().UTC().Unix()
	affected, err = c.repository.UpdateDetail(ctx, req, updatedTime)

	return
}

func (c *CustomerUsecase) ChangeStatus(ctx context.Context, req *pb.CustomerChangeStatusRequest) (affected bool, err error) {
	updatedTime := time.Now().UTC().Unix()
	affected, err = c.repository.ChangeStatus(ctx, req, updatedTime)

	return
}

func (c *CustomerUsecase) FindAll(ctx context.Context, req *pb.CustomerFindAllRequest) (customers *pb.CustomerFindAllResponse, err error) {
	customers, err = c.repository.FindAll(ctx, req)

	return
}

func (c *CustomerUsecase) Find(ctx context.Context, req *pb.CustomerFindRequest) (res *pb.CustomerFindResponse, err error) {
	customer, err := c.repository.Find(ctx, req)
	if err == mongo.ErrNoDocuments {
		res = &pb.CustomerFindResponse{
			IsEmpty: true,
		}
		err = nil

		return
	}

	res = &pb.CustomerFindResponse{
		Payload: customer,
	}

	return
}

func (c *CustomerUsecase) Save(ctx context.Context, req *pb.CustomerCreateRequest) (err error) {
	t := time.Now()
	createdTime := t.Unix()
	generatedId := strconv.Itoa(int(t.UnixNano()))

	err = c.repository.Save(ctx, req, generatedId, createdTime)

	return
}
