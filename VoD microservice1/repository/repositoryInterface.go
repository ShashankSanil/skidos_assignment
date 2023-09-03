package repository

import (
	"VoD_microservice1/common"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RepositortInterface interface {
	Upload(videos []interface{}, language string) common.HTTPResponse
	GetAll(filterObj map[string]interface{}, sortOrder map[string]interface{}, page, size int, language string) common.HTTPResponse
	GetById(id primitive.ObjectID, language string) common.HTTPResponse
}
