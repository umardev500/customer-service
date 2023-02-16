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

func (c *CustomerRepository) Login(ctx context.Context, req *pb.CustomerLoginRequest) (res *pb.CustomerLoginResponse, err error) {
	user := req.User
	pass := req.Pass
	filter := bson.M{"user": user, "pass": pass}

	var result domain.CustomerLoginResponse
	err = c.customers.FindOne(ctx, filter).Decode(&result)
	if err == mongo.ErrNoDocuments {
		res = &pb.CustomerLoginResponse{
			IsEmpty: true,
		}
		err = nil

		return
	}

	res = &pb.CustomerLoginResponse{
		Payload: &pb.CustomerLoginPayload{
			CustomerId: result.CustomerId,
			User:       result.User,
		},
	}

	return
}

func (c *CustomerRepository) SetExp(ctx context.Context, req *pb.CustomerSetExpRequest, updatedTime int64) (affected bool, err error) {
	filter := bson.M{"customer_id": req.CustomerId}
	payload := bson.M{
		"exp_until":       req.ExpTime,
		"status":          req.Status,
		"settlement_time": req.SettlementTime,
		"updated_at":      updatedTime,
	}
	set := bson.M{"$set": payload}
	resp, err := c.customers.UpdateOne(ctx, filter, set)
	if err != nil {
		return false, err
	}

	if resp.ModifiedCount > 0 {
		affected = true
		err = nil
	}

	return
}

func (c *CustomerRepository) Delete(ctx context.Context, req *pb.CustomerDeleteRequest, deletedTime int64) (affected bool, err error) {
	filter := bson.M{"customer_id": req.CustomerId}

	if !req.Hard {
		payload := bson.M{"deleted_at": deletedTime}
		set := bson.M{"$set": payload}
		resp, err := c.customers.UpdateOne(ctx, filter, set)
		if resp.ModifiedCount > 0 {
			return true, err
		}
	}

	resp, err := c.customers.DeleteOne(ctx, filter)
	if resp.DeletedCount > 0 {
		affected = true
	}

	return
}

func (c *CustomerRepository) UpdateDetail(ctx context.Context, req *pb.CustomerUpdateDetailRequest, updatedTime int64) (affected bool, err error) {
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
	resp, err := c.customers.UpdateOne(ctx, filter, set)
	if err != nil {
		return
	}

	// Check if the number of documents modified by the update operation is greater than 0
	if resp.ModifiedCount > 0 {
		// Set affected to true if there were any modified documents
		affected = true
	}

	return
}

func (c *CustomerRepository) ChangeStatus(ctx context.Context, req *pb.CustomerChangeStatusRequest, updatedTime int64) (affected bool, err error) {
	filter := bson.M{"customer_id": req.CustomerId}
	payload := bson.M{
		"status":     req.Status,
		"updated_at": updatedTime,
	}
	set := bson.M{"$set": payload}
	resp, err := c.customers.UpdateOne(ctx, filter, set)

	if resp.ModifiedCount > 0 {
		affected = true
	}

	return
}

func (c *CustomerRepository) FindAll(ctx context.Context, req *pb.CustomerFindAllRequest) (customers *pb.CustomerFindAllResponse, err error) {
	s := req.Search
	status := bson.M{"status": req.Status}
	isExpired := bson.M{}
	if req.Status == "" || req.Status == "none" || req.Status == "deleted" {
		status = bson.M{"status": bson.M{"$ne": nil}}
	}

	deleted := bson.M{"deleted_at": bson.M{"$eq": nil}}

	if req.Status == "deleted" {
		deleted = bson.M{"deleted_at": bson.M{"$ne": nil}}
	}

	if req.Status == "expired" {
		now := time.Now().UTC().Unix()
		status = bson.M{} // reset status
		isExpired = bson.M{"exp_until": bson.M{"$lte": now}}
	}

	today := bson.M{}
	if req.Status == "today" {
		now := time.Now().UTC()
		startOfToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
		startOfTomorrow := startOfToday.Add(24 * time.Hour)

		status = bson.M{}
		today = bson.M{"created_at": bson.M{"$gte": startOfToday.Unix(), "$lte": startOfTomorrow.Unix()}}
	}

	filterData := []bson.M{
		status,
		deleted,
		isExpired,
		today,
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
			{
				"detail.npsn": bson.M{
					"$regex": primitive.Regex{
						Pattern: s,
						Options: "i",
					},
				},
			},
			{
				"user": bson.M{
					"$regex": primitive.Regex{
						Pattern: s,
						Options: "i",
					},
				},
			},
			{
				"status": bson.M{
					"$regex": primitive.Regex{
						Pattern: s,
						Options: "i",
					},
				},
			},
			{
				"detail.name": bson.M{
					"$regex": primitive.Regex{
						Pattern: s,
						Options: "i",
					},
				},
			},
			{
				"detail.email": bson.M{
					"$regex": primitive.Regex{
						Pattern: s,
						Options: "i",
					},
				},
			},
			{
				"detail.wa": bson.M{
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
		cur, err := c.customers.Find(ctx, filter, findOpt)
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

			customer := c.parseCustomerResponse(each)

			customers.Customers = append(customers.Customers, customer)
		}
	}

	rows, _ := c.customers.CountDocuments(ctx, filter)

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

func (c *CustomerRepository) Find(ctx context.Context, req *pb.CustomerFindRequest) (customer *pb.Customer, err error) {
	var data domain.Customer

	filter := bson.M{"customer_id": req.CustomerId}

	err = c.customers.FindOne(ctx, filter).Decode(&data)
	if err != nil {
		return
	}

	customer = c.parseCustomerResponse(data)

	return
}

func (c *CustomerRepository) Save(ctx context.Context, req *pb.CustomerCreateRequest, generatedId string, createdTime int64) (err error) {
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

	_, err = c.customers.InsertOne(ctx, payload)

	return
}
