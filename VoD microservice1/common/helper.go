package common

import (
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chassis/openlog"
)

type SignedDetails struct {
	Email     string
	Username  string
	Uid       string
	User_type string
	jwt.StandardClaims
}

func ErrorHandler(errcode string, response interface{}, totalcount int64, language string) HTTPResponse {
	if language == "" {
		language = "en"
	}

	ErrDetailes := Errorcodes[errcode].(map[string]interface{})
	msg := ErrDetailes[language].(string)
	Status := ErrDetailes["status"].(string)
	status, err := strconv.ParseInt(Status, 10, 64)
	if err != nil {
		openlog.Error(err.Error())
		return HTTPResponse{Msg: "Internal server error", Status: 500}
	}
	code, err := strconv.ParseInt(errcode, 10, 64)
	if err != nil {
		openlog.Error(err.Error())
		return HTTPResponse{Msg: "Internal server error", Status: 500}
	}
	return HTTPResponse{Msg: msg, Status: int(status), Data: response, TotalCount: totalcount, ErrorCode: code}
}
