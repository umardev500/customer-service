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
		Logo:     data.Detail.Logo,
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
			Name:       result.Detail.Name,
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

	location := bson.D{}
	if req.Detail.Location != nil {
		location = bson.D{
			{Key: "detail.location.address", Value: req.Detail.Location.Address},
			{Key: "detail.location.village", Value: req.Detail.Location.Village},
			{Key: "detail.location.district", Value: req.Detail.Location.District},
			{Key: "detail.location.city", Value: req.Detail.Location.City},
			{Key: "detail.location.province", Value: req.Detail.Location.Province},
			{Key: "detail.location.postal_code", Value: req.Detail.Location.PostalCode},
		}

		helper.NoEmpty(location, &location)
	}

	detail := bson.D{
		{Key: "detail.npsn", Value: req.Detail.Npsn},
		{Key: "detail.name", Value: req.Detail.Name},
		{Key: "detail.email", Value: req.Detail.Email},
		{Key: "detail.wa", Value: req.Detail.Wa},
		{Key: "detail.type", Value: req.Detail.Type},
		{Key: "detail.level", Value: req.Detail.Level},
		{Key: "detail.about", Value: req.Detail.About},
		{Key: "detail.logo", Value: req.Detail.Logo},
	}
	helper.NoEmpty(detail, &detail)

	payload := bson.D{}
	payload = append(payload, detail...)
	payload = append(payload, location...)
	payload = append(payload, bson.E{Key: "updated_at", Value: updatedTime})

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

func (c *CustomerRepository) FindAll(ctx context.Context, req *pb.CustomerFindAllRequest) (result *pb.CustomerFindAllResponse, err error) {
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

	result = &pb.CustomerFindAllResponse{}

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

	var customers = []*pb.Customer{}
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

			customers = append(customers, customer)
		}
	}

	result.Payload = &pb.CustomerFindAllPayload{
		Customers: customers,
	}

	rows, _ := c.customers.CountDocuments(ctx, filter)
	if rows < 1 {
		result.IsEmpty = true
	}

	dataSize := int64(len(customers))
	result.Payload.Rows = rows
	result.Payload.Pages = int64(math.Ceil(float64(rows) / float64(perPage)))
	if dataSize < 1 {
		result.Payload.Pages = 0
	} else if perPage == 0 {
		result.Payload.Pages = 1
	}

	result.Payload.PerPage = perPage
	result.Payload.ActivePage = page + 1
	if dataSize < 1 {
		result.Payload.ActivePage = 0
	}
	result.Payload.Total = dataSize

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
