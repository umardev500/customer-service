package injector

import (
	"customer/app/delivery"
	"customer/app/repository"
	"customer/app/usecase"

	"go.mongodb.org/mongo-driver/mongo"
)

func NewCustomerInjector(db *mongo.Database) *delivery.CustomerDelivery {
	repo := repository.NewCustomerRepository(db)
	usecase := usecase.NewCustomerUsecase(repo)

	return delivery.NewCustomerDelivery(usecase)
}
