package common

import "github.com/dgrijalva/jwt-go"

var DbName string

var CollectionName = "user"

var SECRET_KEY string

var Errorcodes = make(map[string]interface{})

type HTTPResponse struct {
	Msg        string      `json:"_msg"`
	Status     int         `json:"_status"`
	Data       interface{} `json:"_data"`
	TotalCount int64       `json:"_totalcount"`
	ErrorCode  int64       `json:"_errorcode"`
}

type SignedDetails struct {
	Email     string
	Username  string
	Uid       string
	User_type string
	jwt.StandardClaims
}
