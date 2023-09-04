package repository

import (
	"VoD_microservice1/common"
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository struct {
	DBClient *mongo.Client
	DBName   string
}

func (tr *Repository) Upload(videos []interface{}, language string) common.HTTPResponse {
	collection := tr.DBClient.Database(tr.DBName).Collection(common.CollectionName)
	result, err := collection.InsertMany(context.TODO(), videos)
	if err != nil {
		return common.ErrorHandler("104", nil, 0, language)
	}
	return common.ErrorHandler("105", result, 0, language)
}

func (tr *Repository) GetAll(filters map[string]interface{}, sortOrder map[string]interface{}, page, size int, language string) common.HTTPResponse {
	collection := tr.DBClient.Database(tr.DBName).Collection(common.CollectionName)
	startIndex := (page - 1) * size
	log.Println(startIndex)
	options := *options.Find()
	result := make([]map[string]interface{}, 0)
	count, err := collection.CountDocuments(context.TODO(), filters)
	if err != nil {
		fmt.Println(err.Error())
		return common.ErrorHandler("106", nil, 0, language)
	}
	cur, err := collection.Find(context.TODO(), filters, options.SetLimit(int64(size)), options.SetSkip((int64(page)-1)*int64(size)), options.SetSort(sortOrder))
	if err != nil {
		fmt.Println(err.Error())
		return common.ErrorHandler("106", nil, 0, language)
	}
	for cur.Next(context.TODO()) {
		var doc map[string]interface{}
		err := cur.Decode(&doc)
		if err != nil {
			return common.ErrorHandler("106", nil, 0, language)
		}
		result = append(result, doc)
	}
	return common.ErrorHandler("107", result, count, language)
}

func (tr *Repository) GetById(id primitive.ObjectID, language string) common.HTTPResponse {
	collection := tr.DBClient.Database(tr.DBName).Collection(common.CollectionName)
	result := make(map[string]interface{})
	err := collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return common.ErrorHandler("112", nil, 0, language)
		}
		return common.ErrorHandler("108", nil, 0, language)
	}
	return common.ErrorHandler("109", result, 1, language)
}
