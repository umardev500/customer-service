package usecase

import (
	"context"
	"customer/pb"
)

func (c *CustomerUsecase) UpdateCreds(ctx context.Context, req *pb.CustomerUpdateCredsRequest) (res *pb.OperationResponse, err error) {
	return c.repository.UpdateCreds(ctx, req)
}
