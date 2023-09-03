package service

import (
	"VoD_microservice/common"
	repository "VoD_microservice/repository"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/go-chassis/go-chassis/v2/server/restful"
)

type Service struct {
	Rep repository.RepositortInterface
}

func (s *Service) Signup(user map[string]interface{}, language string) common.HTTPResponse {
	httres := s.Rep.Signup(user, language)
	return httres

}

func (s *Service) Login(user map[string]interface{}, language string) common.HTTPResponse {
	httres := s.Rep.Login(user, language)
	return httres

}

func (s *Service) GetAllUser(context *restful.Context, language string) common.HTTPResponse {
	recordPerPage, err := strconv.Atoi(context.Req.Request.URL.Query().Get("size"))
	if err != nil || recordPerPage < 1 {
		recordPerPage = 10
	}
	page, err1 := strconv.Atoi(context.Req.Request.URL.Query().Get("page"))
	if err1 != nil || page < 1 {
		page = 1
	}
	filter := context.Req.Request.URL.Query().Get("filter")
	filterObj := make(map[string]interface{})
	byte, err2 := json.Marshal(filter)
	if err2 != nil {
		fmt.Println(err2.Error())
		common.ErrorHandler("118", nil, 0, language)
	}
	err3 := json.Unmarshal(byte, &filterObj)
	if err3 != nil {
		fmt.Println(err3.Error())
		common.ErrorHandler("118", nil, 0, language)
	}
	return s.Rep.GetUserByPagination(filterObj, map[string]interface{}{"created_at": -1}, page, recordPerPage, language)
}
