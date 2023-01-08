package repository

import (
	"context"
	"customer/domain"
	"customer/helper"
	"customer/pb"
	"math"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (pr *CustomerRepository) parseCustomerResponse(data domain.Customer) (customer *pb.Customer) {
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

// Template
// func (pr *CustomerRepository) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()
// }

func (pr *CustomerRepository) Delete(req *pb.CustomerDeleteRequest, deletedTime int64) (affected bool, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"customer_id": req.CustomerId}

	if !req.Hard {
		payload := bson.M{"deleted_at": deletedTime}
		set := bson.M{"$set": payload}
		resp, err := pr.customers.UpdateOne(ctx, filter, set)
		if resp.ModifiedCount > 0 {
			return true, err
		}
	}

	resp, err := pr.customers.DeleteOne(ctx, filter)
	if resp.DeletedCount > 0 {
		affected = true
	}

	return
}

func (pr *CustomerRepository) UpdateDetail(req *pb.CustomerUpdateDetailRequest, updatedTime int64) (affected bool, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"customer_id": req.CustomerId}

	var location bson.M

	if req.Detail.Location != nil {
		location = bson.M{
			"address":     req.Detail.Location.Address,
			"village":     req.Detail.Location.Village,
			"district":    req.Detail.Location.District,
			"city":        req.Detail.Location.City,
			"province":    req.Detail.Location.Province,
			"postal_code": req.Detail.Location.PostalCode,
		}
	}

	detail := bson.D{
		{Key: "npsn", Value: req.Detail.Npsn},
		{Key: "name", Value: req.Detail.Name},
		{Key: "email", Value: req.Detail.Email},
		{Key: "wa", Value: req.Detail.Wa},
		{Key: "type", Value: req.Detail.Type},
		{Key: "level", Value: req.Detail.Level},
		{Key: "about", Value: req.Detail.About},
		{Key: "location", Value: location},
	}

	// Initialize a slice of bson.Elements called detailFix
	detailFix := bson.D{}

	// Call the NoEmpty function on the detail variable and pass in a pointer to detailFix
	helper.NoEmpty(detail, &detailFix)

	payload := bson.M{
		"detail": detailFix,
	}
	set := bson.M{"$set": payload}
	resp, err := pr.customers.UpdateOne(ctx, filter, set)

	// Check if the number of documents modified by the update operation is greater than 0
	if resp.ModifiedCount > 0 {
		// Set affected to true if there were any modified documents
		affected = true
	}

	return
}

func (pr *CustomerRepository) ChangeStatus(req *pb.CustomerChangeStatusRequest, updatedTime int64) (affected bool, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"customer_id": req.CustomerId}
	payload := bson.M{
		"status":     req.Status,
		"updated_at": updatedTime,
	}
	set := bson.M{"$set": payload}
	resp, err := pr.customers.UpdateOne(ctx, filter, set)

	if resp.ModifiedCount > 0 {
		affected = true
	}

	return
}

func (pr *CustomerRepository) FindAll(req *pb.CustomerFindAllRequest) (customers *pb.CustomerFindAllResponse, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	s := req.Search
	status := bson.M{"status": req.Status}
	if (req.Status == "" || req.Status == "none") && req.Status != "deleted" {
		status = bson.M{"status": bson.M{"$ne": nil}}
	}

	deleted := bson.M{"deleted_at": bson.M{"$eq": nil}}

	filterData := []bson.M{
		status,
		deleted,
	}

	if req.Status == "deleted" {
		deleted = bson.M{"deleted_at": bson.M{"$ne": nil}}
		filterData = []bson.M{deleted}
	}

	customers = &pb.CustomerFindAllResponse{}

	filter := bson.M{
		"$or": []bson.M{
			{
				"customer_id": bson.M{
					"$regex": primitive.Regex{
						Pattern: s,
						Options: "i",
					},
				},
			},
		},
		"$and": filterData,
	}

	findOpt := options.Find()

	if req.Sort == "desc" {
		findOpt.SetSort(bson.M{"customer_id": -1})
	}

	page := req.Page
	perPage := req.PerPage
	offset := page * perPage

	findOpt.SetSkip(offset)
	findOpt.SetLimit(perPage)

	if !req.CountOnly {
		cur, err := pr.customers.Find(ctx, filter, findOpt)
		if err != nil {
			return nil, err
		}

		defer cur.Close(ctx)

		for cur.Next(ctx) {
			var each domain.Customer

			err = cur.Decode(&each)
			if err != nil {
				return nil, err
			}

			customer := pr.parseCustomerResponse(each)
			customers.Customers = append(customers.Customers, customer)
		}
	}

	rows, _ := pr.customers.CountDocuments(ctx, filter)

	dataSize := int64(len(customers.Customers))
	customers.Rows = rows
	customers.Pages = int64(math.Ceil(float64(rows) / float64(perPage)))
	if dataSize < 1 {
		customers.Pages = 0
	} else if perPage == 0 {
		customers.Pages = 1
	}

	customers.PerPage = perPage
	customers.ActivePage = page + 1
	if dataSize < 1 {
		customers.ActivePage = 0
	}
	customers.Total = dataSize

	return
}

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

	customer = pr.parseCustomerResponse(data)

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
