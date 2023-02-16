package usecase

import (
	"context"
	"customer/domain"
	"customer/pb"
	"strconv"
	"time"
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

func (pu *CustomerUsecase) SetExp(ctx context.Context, req *pb.CustomerSetExpRequest) (affected bool, err error) {
	updatedTime := time.Now().UTC().Unix()
	affected, err = pu.repository.SetExp(ctx, req, updatedTime)

	return
}

func (pu *CustomerUsecase) Delete(ctx context.Context, req *pb.CustomerDeleteRequest) (affected bool, err error) {
	deletedTime := time.Now().UTC().Unix()
	affected, err = pu.repository.Delete(ctx, req, deletedTime)

	return
}

func (pu *CustomerUsecase) UpdateDetail(ctx context.Context, req *pb.CustomerUpdateDetailRequest) (affected bool, err error) {
	updatedTime := time.Now().UTC().Unix()
	affected, err = pu.repository.UpdateDetail(ctx, req, updatedTime)

	return
}

func (pu *CustomerUsecase) ChangeStatus(ctx context.Context, req *pb.CustomerChangeStatusRequest) (affected bool, err error) {
	updatedTime := time.Now().UTC().Unix()
	affected, err = pu.repository.ChangeStatus(ctx, req, updatedTime)

	return
}

func (pu *CustomerUsecase) FindAll(ctx context.Context, req *pb.CustomerFindAllRequest) (customers *pb.CustomerFindAllResponse, err error) {
	customers, err = pu.repository.FindAll(ctx, req)

	return
}

func (pu *CustomerUsecase) FindOne(ctx context.Context, req *pb.CustomerFindOneRequest) (customer *pb.Customer, err error) {
	customer, err = pu.repository.FindOne(ctx, req)

	return
}

func (pu *CustomerUsecase) Save(ctx context.Context, req *pb.CustomerCreateRequest) (err error) {
	t := time.Now()
	createdTime := t.Unix()
	generatedId := strconv.Itoa(int(t.UnixNano()))

	err = pu.repository.Save(ctx, req, generatedId, createdTime)

	return
}
