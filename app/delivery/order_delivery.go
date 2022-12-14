package delivery

import (
	"context"
	"customer/domain"
	"customer/pb"
)

type CustomerDelivery struct {
	usecase domain.CustomerUsecase
	pb.UnimplementedCustomerServiceServer
}

func NewCustomerDelivery(usecase domain.CustomerUsecase) *CustomerDelivery {
	return &CustomerDelivery{
		usecase: usecase,
	}
}

// Template
//func (pd *CustomerDelivery) Delete(ctx context.Context, req *pb.) (res *pb., err error) {}

func (pd *CustomerDelivery) FindOne(ctx context.Context, req *pb.CustomerFindOneRequest) (customer *pb.Customer, err error) {
	customer, err = pd.usecase.FindOne(req)

	return
}

func (pd *CustomerDelivery) Create(ctx context.Context, req *pb.CustomerCreateRequest) (res *pb.Empty, err error) {
	res = &pb.Empty{}
	err = pd.usecase.Save(req)

	return
}
