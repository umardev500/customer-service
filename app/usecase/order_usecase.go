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
