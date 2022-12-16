package domain

import "customer/pb"

type CustomerUsecase interface {
	Save(req *pb.CustomerCreateRequest) error
}

type CustomerRepository interface {
	Save(req *pb.CustomerCreateRequest, generatedId string, createdTime int64) error
}
