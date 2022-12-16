package repository

import (
	"context"
	"customer/domain"
	"customer/pb"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type CustomerRepository struct {
	db        *mongo.Database
	customers *mongo.Collection
}

func NewCustomerRepository(db *mongo.Database) domain.CustomerRepository {
	return &CustomerRepository{
		db:        db,
		customers: db.Collection("customers"),
	}
}

// Template
// func (pr *CustomerRepository) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()
// }

func (pr *CustomerRepository) Save(req *pb.CustomerCreateRequest, generatedId string, createdTime int64) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	status := "pending"

	detail := bson.D{
		{Key: "name", Value: req.Detail.Name},
		{Key: "email", Value: req.Detail.Email},
		{Key: "wa", Value: req.Detail.Wa},
	}

	payload := bson.D{
		{Key: "customer_id", Value: generatedId},
		{Key: "user", Value: req.User},
		{Key: "pass", Value: req.Pass},
		{Key: "detail", Value: detail},
		{Key: "status", Value: status},
		{Key: "created_at", Value: createdTime},
	}

	_, err = pr.customers.InsertOne(ctx, payload)

	return
}
