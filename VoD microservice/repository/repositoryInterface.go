package repository

import (
	"VoD_microservice/common"
)

type RepositortInterface interface {
	Signup(user map[string]interface{}, language string) common.HTTPResponse
	Login(user map[string]interface{}, language string) common.HTTPResponse
	GetUserByPagination(filter map[string]interface{}, sortorder map[string]interface{}, page, size int, language string) common.HTTPResponse
}
