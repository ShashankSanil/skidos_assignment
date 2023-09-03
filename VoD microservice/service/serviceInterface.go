package service

import (
	"VoD_microservice/common"
	"github.com/go-chassis/go-chassis/v2/server/restful"
)

type ServiceInterface interface {
	Signup(user map[string]interface{}, language string) common.HTTPResponse
	Login(user map[string]interface{}, language string) common.HTTPResponse
	GetAllUser(context *restful.Context, language string) common.HTTPResponse
}
