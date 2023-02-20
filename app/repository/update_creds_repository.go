package repository

import (
	"context"
	"customer/pb"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func (c *CustomerRepository) UpdateCreds(ctx context.Context, req *pb.CustomerUpdateCredsRequest) (res *pb.OperationResponse, err error) {
	affected := false
	now := time.Now().UTC().Unix()

	filter := bson.M{"user": req.User, "pass": req.Pass}
	payload := bson.M{
		"pass":       req.NewPass,
		"updated_at": now,
	}
	set := bson.M{"$set": payload}
	resp, err := c.customers.UpdateOne(ctx, filter, set)
	if err != nil {
		return
	}

	if resp.ModifiedCount > 0 {
		affected = true
	}

	res = &pb.OperationResponse{IsAffected: affected}

	return
}
