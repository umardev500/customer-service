package helper

import (
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
)

func NoEmpty(source bson.D, dest *bson.D) {
	out := bson.D{}

	for _, val := range source {
		isZero := reflect.ValueOf(val.Value).IsZero()

		if !isZero {
			each := bson.E{Key: val.Key, Value: val.Value}
			out = append(out, each)
		}
	}

	*dest = out
}
