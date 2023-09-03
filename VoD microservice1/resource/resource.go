package resource

import (
	"VoD_microservice1/common"
	service "VoD_microservice1/service"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/go-chassis/go-chassis/v2/pkg/tool"
	"github.com/go-chassis/go-chassis/v2/server/restful"
	"github.com/go-chassis/openlog"
)

type Resources struct {
	ServiceProvider func(dbname string) *service.Service
}

func WriterResponse(context *restful.Context, errcode string, language string) {
	ErrDetails := common.Errorcodes[errcode].(map[string]interface{})
	Msg := ErrDetails[language].(string)
	Status := ErrDetails["status"].(string)
	status, err := strconv.ParseInt(Status, 10, 64)
	if err != nil {
		openlog.Error(err.Error())
		WriterResponse(context, errcode, language)
		return
	}
	code, err := strconv.ParseInt(errcode, 10, 64)
	if err != nil {
		openlog.Error(err.Error())
		WriterResponse(context, errcode, language)
		return
	}
	context.WriteHeaderAndJSON(400, common.HTTPResponse{Msg: Msg, Status: int(status), ErrorCode: code}, "application/json")
}

func Recover(context *restful.Context) {
	if r := recover(); r != nil {
		var stacktrace = tool.GetStackTrace(3)
		openlog.Error("handle request panic.", openlog.WithTags(openlog.Tags{
			"panic": r,
			"stack": stacktrace,
		}))
		res := common.HTTPResponse{Msg: "Invalid Error", Status: 500, TotalCount: 0}
		openlog.Debug("panic recovered")
		if err := context.WriteHeaderAndJSON(res.Status, res, "application/json"); err != nil {
			openlog.Error("write response failed when handler panic.", openlog.WithTags(openlog.Tags{
				"err": err.Error(),
			}))
		}
	}
}

func (s *Resources) Upload(context *restful.Context) {
	openlog.Info("Got a request for Upload")
	defer Recover(context)
	language := context.Req.Request.Header.Get("language")
	if language == "" {
		WriterResponse(context, "101", "en")
		return
	}
	videos := make([]interface{}, 0)
	errbytes, err := ioutil.ReadFile("./conf/videos.json")
	if err != nil {
		WriterResponse(context, "102", "en")
		return
	}
	err = json.Unmarshal(errbytes, &videos)
	if err != nil {
		WriterResponse(context, "103", "en")
		return
	}
	ts := s.ServiceProvider(common.DbName)
	httpres := ts.Upload(videos, language)
	context.WriteHeaderAndJSON(httpres.Status, httpres, "application/json")
}

func (s *Resources) GetAllVideos(context *restful.Context) {
	openlog.Info("Got a request for GetAllVideos")
	defer Recover(context)
	language := context.Req.Request.Header.Get("language")
	if language == "" {
		WriterResponse(context, "101", "en")
		return
	}
	size, err := strconv.Atoi(context.Req.Request.URL.Query().Get("size"))
	if err != nil || size < 1 {
		size = 10
	}
	page, err1 := strconv.Atoi(context.Req.Request.URL.Query().Get("page"))
	if err1 != nil || page < 1 {
		page = 1
	}
	filter := context.Req.Request.URL.Query().Get("filter")
	filterObj := make(map[string]interface{})
	err3 := json.Unmarshal([]byte(filter), &filterObj)
	if err3 != nil {
		fmt.Println(err3.Error())
		common.ErrorHandler("103", nil, 0, language)
	}
	ts := s.ServiceProvider(common.DbName)
	httpres := ts.GetAllVideos(filterObj, page, size, language)
	context.WriteHeaderAndJSON(httpres.Status, httpres, "application/json")
}

func (s *Resources) GetVideoById(context *restful.Context) {
	openlog.Info("Got a request for GetVideoById")
	defer Recover(context)
	language := context.Req.Request.Header.Get("language")
	if language == "" {
		WriterResponse(context, "101", "en")
		return
	}
	id := context.ReadPathParameter("id")
	ts := s.ServiceProvider(common.DbName)
	httpres := ts.GetVideoById(id, language)
	context.WriteHeaderAndJSON(httpres.Status, httpres, "application/json")
}

func (s *Resources) URLPatterns() []restful.Route {
	return []restful.Route{
		//video routes
		{Method: http.MethodPost, Path: "/video/upload", ResourceFunc: s.Upload},
		{Method: http.MethodGet, Path: "/video/{vid}", ResourceFunc: s.GetVideoById},
		{Method: http.MethodGet, Path: "/videos", ResourceFunc: s.GetAllVideos},
	}

}
