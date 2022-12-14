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

func (pr *CustomerRepository) FindOne(req *pb.CustomerFindOneRequest) (customer *pb.Customer, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var data domain.Customer

	filter := bson.M{
		"$or": []bson.M{
			{
				"customer_id": req.CustomerId,
			},
			{
				"$and": []bson.M{
					{"user": req.User},
					{"pass": req.Pass},
				},
			},
		},
	}

	err = pr.customers.FindOne(ctx, filter).Decode(&data)
	if err != nil {
		return
	}

	var location *pb.CustomerLocation

	if data.Detail.Location != nil {
		location = &pb.CustomerLocation{
			Address:    data.Detail.Location.Address,
			Village:    data.Detail.Location.Village,
			District:   data.Detail.Location.District,
			City:       data.Detail.Location.City,
			Province:   data.Detail.Location.Province,
			PostalCode: data.Detail.Location.PostalCode,
		}
	}

	detail := &pb.CustomerDetail{
		Npsn:     data.Detail.Npsn,
		Name:     data.Detail.Name,
		Email:    data.Detail.Email,
		Wa:       data.Detail.Wa,
		Type:     data.Detail.Type,
		Level:    data.Detail.Level,
		About:    data.Detail.About,
		Location: location,
	}

	customer = &pb.Customer{
		CustomerId: data.CustomerId,
		User:       data.User,
		Pass:       data.Pass,
		Detail:     detail,
		ExpUntil:   data.ExpUntil,
		Status:     data.Status,
		CreatedAt:  data.CreatedAt,
		UpdatedAt:  data.UpdatedAt,
		DeletedAt:  data.DeletedAt,
	}

	return
}

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
