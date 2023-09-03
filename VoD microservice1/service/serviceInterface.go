package service

import "VoD_microservice1/common"

type ServiceInterface interface {
	Upload(videos []interface{}, language string) common.HTTPResponse
	GetAllVideos(filterObj map[string]interface{}, page, size int, language string) common.HTTPResponse
	GetVideoById(id, language string) common.HTTPResponse
}
