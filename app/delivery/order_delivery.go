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

func (pd *CustomerDelivery) ChangeStatus(ctx context.Context, req *pb.CustomerChangeStatusRequest) (res *pb.OperationResponse, err error) {
	affected, err := pd.usecase.ChangeStatus(req)
	res = &pb.OperationResponse{IsAffected: affected}

	return
}

func (pd *CustomerDelivery) FindAll(ctx context.Context, req *pb.CustomerFindAllRequest) (customer *pb.CustomerFindAllResponse, err error) {
	customer, err = pd.usecase.FindAll(req)

	return
}

func (pd *CustomerDelivery) FindOne(ctx context.Context, req *pb.CustomerFindOneRequest) (customer *pb.Customer, err error) {
	customer, err = pd.usecase.FindOne(req)

	return
}

func (pd *CustomerDelivery) Create(ctx context.Context, req *pb.CustomerCreateRequest) (res *pb.Empty, err error) {
	res = &pb.Empty{}
	err = pd.usecase.Save(req)

	return
}
