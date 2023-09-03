package service

import (
	"VoD_microservice1/common"
	repository "VoD_microservice1/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	Rep repository.RepositortInterface
}

func (s *Service) Upload(videos []interface{}, language string) common.HTTPResponse {
	return s.Rep.Upload(videos, language)
}

func (s *Service) GetAllVideos(filterObj map[string]interface{}, page, size int, language string) common.HTTPResponse {
	sortOrder := map[string]interface{}{"title": -1}
	return s.Rep.GetAll(filterObj, sortOrder, page, size, language)
}

func (s *Service) GetVideoById(id, language string) common.HTTPResponse {
	obId, _ := primitive.ObjectIDFromHex(id)
	return s.Rep.GetById(obId, language)
}
