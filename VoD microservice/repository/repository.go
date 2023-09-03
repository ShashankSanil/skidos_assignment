package repository

import (
	"VoD_microservice/common"
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository struct {
	DBClient *mongo.Client
	DBName   string
}

// Signup
func (tr *Repository) Signup(user map[string]interface{}, language string) common.HTTPResponse {
	collection := tr.DBClient.Database(tr.DBName).Collection(common.CollectionName)

	count, err := collection.CountDocuments(context.TODO(), bson.M{"Email": user["Email"].(string)})
	if err != nil {
		log.Panic(err)
		return common.ErrorHandler("101", nil, 0, language)
	}
	if count > 0 {
		return common.ErrorHandler("102", nil, 0, language)
	}
	password := common.HashPassword(user["Password"].(string))
	user["Password"] = &password
	user["Created_at"], _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user["Updated_at"], _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user["ID"] = primitive.NewObjectID()
	user["User_id"] = user["ID"].(primitive.ObjectID).Hex()
	delete(user, "ID")
	token, refreshToken, err, errcode := common.GenerateAllTokens(user["Email"].(string), user["Username"].(string), user["User_type"].(string), user["User_id"].(string))
	if err != nil {
		return common.ErrorHandler(errcode, nil, 0, language)
	}
	user["Token"] = &token
	user["Refresh_token"] = &refreshToken

	resp, insertErr := collection.InsertOne(context.TODO(), user)
	if insertErr != nil {
		fmt.Println("User not created !!!")
		return common.ErrorHandler("104", nil, 0, language)
	}
	return common.ErrorHandler("105", resp, 0, language)
}

// Login
func (tr *Repository) Login(user map[string]interface{}, language string) common.HTTPResponse {
	collection := tr.DBClient.Database(tr.DBName).Collection(common.CollectionName)

	var foundUser map[string]interface{}
	err := collection.FindOne(context.TODO(), bson.M{"Email": user["Email"].(string)}).Decode(&foundUser)
	if err != nil {
		return common.ErrorHandler("108", nil, 0, language)
	}

	passwordIsValid, errcode := common.VerifyPassword(user["Password"].(string), foundUser["Password"].(string))
	if !passwordIsValid {
		return common.ErrorHandler(errcode, nil, 0, language)
	}

	if foundUser["Email"] == nil {
		return common.ErrorHandler("110", nil, 0, language)
	}

	token, refreshToken, err, errcode := common.GenerateAllTokens(foundUser["Email"].(string), foundUser["Username"].(string), foundUser["User_type"].(string), foundUser["User_id"].(string))
	if err != nil {
		return common.ErrorHandler(errcode, nil, 0, language)
	}
	common.UpdateAllTokens(token, refreshToken, foundUser["User_id"].(string), collection)
	err = collection.FindOne(context.TODO(), bson.M{"User_id": foundUser["User_id"].(string)}).Decode(&foundUser)
	if err != nil {
		return common.ErrorHandler("108", nil, 0, language)
	}
	return common.ErrorHandler("111", foundUser, 1, language)
}

// GetAll
func (tr *Repository) GetUserByPagination(filters map[string]interface{}, sortorder map[string]interface{}, page, size int, language string) common.HTTPResponse {
	collection := tr.DBClient.Database(tr.DBName).Collection(common.CollectionName)
	startIndex := (page - 1) * size
	log.Println(startIndex)
	options := *options.Find()
	result := make([]map[string]interface{}, 0)
	cur, err := collection.Find(context.TODO(), filters, options.SetLimit(int64(size)), options.SetSkip((int64(page)-1)*int64(size)), options.SetSort(sortorder))
	if err != nil {
		return common.ErrorHandler("116", nil, 0, language)
	}
	for cur.Next(context.TODO()) {
		var doc map[string]interface{}
		err := cur.Decode(&doc)
		if err != nil {
			return common.ErrorHandler("116", nil, 0, language)
		}
		result = append(result, doc)
	}
	return common.ErrorHandler("117", result, int64(len(result)), language)
}
