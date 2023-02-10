package usecase

import (
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

func (pu *CustomerUsecase) SetExp(req *pb.CustomerSetExpRequest) (affected bool, err error) {
	updatedTime := time.Now().UTC().Unix()
	affected, err = pu.repository.SetExp(req, updatedTime)

	return
}

func (pu *CustomerUsecase) Delete(req *pb.CustomerDeleteRequest) (affected bool, err error) {
	deletedTime := time.Now().UTC().Unix()
	affected, err = pu.repository.Delete(req, deletedTime)

	return
}

func (pu *CustomerUsecase) UpdateDetail(req *pb.CustomerUpdateDetailRequest) (affected bool, err error) {
	updatedTime := time.Now().UTC().Unix()
	affected, err = pu.repository.UpdateDetail(req, updatedTime)

	return
}

func (pu *CustomerUsecase) ChangeStatus(req *pb.CustomerChangeStatusRequest) (affected bool, err error) {
	updatedTime := time.Now().UTC().Unix()
	affected, err = pu.repository.ChangeStatus(req, updatedTime)

	return
}

func (pu *CustomerUsecase) FindAll(req *pb.CustomerFindAllRequest) (customers *pb.CustomerFindAllResponse, err error) {
	customers, err = pu.repository.FindAll(req)

	return
}

func (pu *CustomerUsecase) FindOne(req *pb.CustomerFindOneRequest) (customer *pb.Customer, err error) {
	customer, err = pu.repository.FindOne(req)

	return
}

func (pu *CustomerUsecase) Save(req *pb.CustomerCreateRequest) (err error) {
	t := time.Now()
	createdTime := t.Unix()
	generatedId := strconv.Itoa(int(t.UnixNano()))

	err = pu.repository.Save(req, generatedId, createdTime)

	return
}
