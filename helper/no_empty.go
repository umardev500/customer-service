package helper

import (
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func NoEmpty(source bson.D, dest *bson.D) {
	for _, val := range source {
		passToField := func() {
			detailValue := bson.E{Key: val.Key, Value: val.Value}
			*dest = append(*dest, detailValue)
		}

		isMap := reflect.TypeOf(val.Value) == reflect.TypeOf(primitive.M{})
		if isMap {
			if len(val.Value.(primitive.M)) > 0 {
				passToField()
			}
		} else if len(val.Value.(string)) > 0 {
			passToField()
		}
	}
}
