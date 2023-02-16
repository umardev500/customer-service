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
// func (pd *CustomerDelivery) Delete(ctx context.Context, req *pb.) (res *pb., err error) {}

func (pd *CustomerDelivery) SetExp(ctx context.Context, req *pb.CustomerSetExpRequest) (res *pb.OperationResponse, err error) {
	affected, err := pd.usecase.SetExp(ctx, req)
	res = &pb.OperationResponse{IsAffected: affected}

	return
}

func (pd *CustomerDelivery) Delete(ctx context.Context, req *pb.CustomerDeleteRequest) (res *pb.OperationResponse, err error) {
	affected, err := pd.usecase.Delete(ctx, req)
	res = &pb.OperationResponse{IsAffected: affected}

	return
}

func (pd *CustomerDelivery) UpdateDetail(ctx context.Context, req *pb.CustomerUpdateDetailRequest) (res *pb.OperationResponse, err error) {
	affected, err := pd.usecase.UpdateDetail(ctx, req)
	res = &pb.OperationResponse{IsAffected: affected}

	return
}

func (pd *CustomerDelivery) ChangeStatus(ctx context.Context, req *pb.CustomerChangeStatusRequest) (res *pb.OperationResponse, err error) {
	affected, err := pd.usecase.ChangeStatus(ctx, req)
	res = &pb.OperationResponse{IsAffected: affected}

	return
}

func (pd *CustomerDelivery) FindAll(ctx context.Context, req *pb.CustomerFindAllRequest) (customer *pb.CustomerFindAllResponse, err error) {
	customer, err = pd.usecase.FindAll(ctx, req)

	return
}

func (pd *CustomerDelivery) FindOne(ctx context.Context, req *pb.CustomerFindOneRequest) (customer *pb.Customer, err error) {
	customer, err = pd.usecase.FindOne(ctx, req)

	return
}

func (pd *CustomerDelivery) Create(ctx context.Context, req *pb.CustomerCreateRequest) (res *pb.Empty, err error) {
	res = &pb.Empty{}
	err = pd.usecase.Save(ctx, req)

	return
}
