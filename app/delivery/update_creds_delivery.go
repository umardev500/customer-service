package delivery

import (
	"context"
	"customer/pb"
)

func (c *CustomerDelivery) UpdateCreds(ctx context.Context, req *pb.CustomerUpdateCredsRequest) (res *pb.OperationResponse, err error) {
	return c.usecase.UpdateCreds(ctx, req)
}
