package delivery

import (
	"context"
	"customer/domain"
	"customer/pb"
	"fmt"
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

func (c *CustomerDelivery) Login(ctx context.Context, req *pb.CustomerLoginRequest) (res *pb.CustomerLoginResponse, err error) {
	res, err = c.usecase.Login(ctx, req)
	return
}

func (c *CustomerDelivery) SetExp(ctx context.Context, req *pb.CustomerSetExpRequest) (res *pb.OperationResponse, err error) {
	affected, err := c.usecase.SetExp(ctx, req)
	res = &pb.OperationResponse{IsAffected: affected}

	return
}

func (c *CustomerDelivery) Delete(ctx context.Context, req *pb.CustomerDeleteRequest) (res *pb.OperationResponse, err error) {
	affected, err := c.usecase.Delete(ctx, req)
	res = &pb.OperationResponse{IsAffected: affected}

	return
}

func (c *CustomerDelivery) UpdateDetail(ctx context.Context, req *pb.CustomerUpdateDetailRequest) (res *pb.OperationResponse, err error) {
	affected, err := c.usecase.UpdateDetail(ctx, req)
	res = &pb.OperationResponse{IsAffected: affected}

	return
}

func (c *CustomerDelivery) ChangeStatus(ctx context.Context, req *pb.CustomerChangeStatusRequest) (res *pb.OperationResponse, err error) {
	affected, err := c.usecase.ChangeStatus(ctx, req)
	res = &pb.OperationResponse{IsAffected: affected}

	return
}

func (c *CustomerDelivery) FindAll(ctx context.Context, req *pb.CustomerFindAllRequest) (customer *pb.CustomerFindAllResponse, err error) {
	customer, err = c.usecase.FindAll(ctx, req)

	return
}

func (c *CustomerDelivery) Find(ctx context.Context, req *pb.CustomerFindRequest) (res *pb.CustomerFindResponse, err error) {
	fmt.Println("find hit")

	res, err = c.usecase.Find(ctx, req)

	return
}

func (c *CustomerDelivery) Create(ctx context.Context, req *pb.CustomerCreateRequest) (res *pb.Empty, err error) {
	res = &pb.Empty{}
	err = c.usecase.Save(ctx, req)

	return
}
