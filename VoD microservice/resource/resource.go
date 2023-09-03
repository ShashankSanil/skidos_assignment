package resource

import (
	"VoD_microservice/common"
	service "VoD_microservice/service"
	"net/http"
	"strconv"

	"github.com/go-chassis/go-chassis/v2/pkg/tool"
	"github.com/go-chassis/go-chassis/v2/server/restful"
	"github.com/go-chassis/openlog"
)

type Resources struct {
	ServiceProvider func(dbname string) *service.Service
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

func (s *Resources) Signup(context *restful.Context) {
	openlog.Info("Got a request for signUp")
	defer Recover(context)
	language := context.Req.Request.Header.Get("language")
	if language == "" {
		WriterResponse(context, "113", "en")
		return
	}
	user := make(map[string]interface{})
	if err := context.ReadEntity(&user); err != nil {
		WriterResponse(context, "106", "en")
		return
	}
	ts := s.ServiceProvider(common.DbName)
	httpres := ts.Signup(user, language)

	context.WriteHeaderAndJSON(httpres.Status, httpres, "application/json")

}

func (s *Resources) Login(context *restful.Context) {
	openlog.Info("Got a request for Login")
	defer Recover(context)
	language := context.Req.Request.Header.Get("language")
	if language == "" {
		WriterResponse(context, "113", "en")
		return
	}
	var user map[string]interface{}
	if err := context.ReadEntity(&user); err != nil {
		WriterResponse(context, "107", "en")
		return
	}

	ts := s.ServiceProvider(common.DbName)
	httpres := ts.Login(user, language)
	context.WriteHeaderAndJSON(httpres.Status, httpres, "application/json")

}
func (s *Resources) GetAllUser(context *restful.Context) {
	openlog.Info("Got a request for GetAllUser")
	defer Recover(context)
	language := context.Req.Request.Header.Get("language")
	if language == "" {
		WriterResponse(context, "113", "en")
		return
	}
	ts := s.ServiceProvider(common.DbName)
	httpres := ts.GetAllUser(context, language)
	context.WriteHeaderAndJSON(httpres.Status, httpres, "application/json")
}

func (s *Resources) URLPatterns() []restful.Route {
	return []restful.Route{
		//user routes
		{Method: http.MethodPost, Path: "/user/signup", ResourceFunc: s.Signup},
		{Method: http.MethodPost, Path: "/user/login", ResourceFunc: s.Login},
		{Method: http.MethodGet, Path: "/users", ResourceFunc: s.GetAllUser},
	}

}
