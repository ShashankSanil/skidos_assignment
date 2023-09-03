package chassisHandlers

import (
	"VoD_microservice1/common"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/emicklei/go-restful"
	"github.com/go-chassis/go-archaius"
	"github.com/go-chassis/go-chassis/v2/core/handler"
	"github.com/go-chassis/go-chassis/v2/core/invocation"
	"github.com/go-chassis/openlog"
)

type Authentication struct{}

func init() { handler.RegisterHandler("Authentication_Handler", AuthenticationHandler) }

func AuthenticationHandler() handler.Handler { return &Authentication{} }

func (h *Authentication) Name() string { return "Authentication_Handler" }

func (h *Authentication) Handle(chain *handler.Chain, inv *invocation.Invocation, cb invocation.ResponseCallBack) {
	openlog.Info("validating token")
	// request object
	var req *http.Request
	if r, ok := inv.Args.(*http.Request); ok {
		req = r
	} else if r, ok := inv.Args.(*restful.Request); ok {
		req = r.Request
	} else {
		openlog.Error(fmt.Sprintf("this handler only works for http protocol, wrong type: %t", inv.Args))
		return
	}
	var resp *http.ResponseWriter
	if r, ok := inv.Reply.(*http.ResponseWriter); ok {
		resp = r
	} else if r, ok := inv.Reply.(*restful.Response); ok {
		resp = &r.ResponseWriter
	} else {
		openlog.Error(fmt.Sprintf("this handler only works for http protocol, wrong type: %t", inv.Args))
		return
	}
	(*resp).Header().Set("Content-Type", "application/json; charset=utf-8")
	if inv.OperationID != "Signup" && inv.OperationID != "Login" {
		token := req.Header.Get("token")
		if token == "" {
			cookie, err := req.Cookie("token")
			if err != nil {
				token = ""
			} else {
				token = string(cookie.Value)
			}
		}
		secretKey := archaius.Get("secretkey").(string)
		tokenDetails, err := jwt.ParseWithClaims(token, &common.SignedDetails{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})
		if err != nil {
			data := common.ErrorHandler("110", nil, 0, "en")
			bytes, _ := json.Marshal(data)
			(*resp).Write(bytes)
			cb(&invocation.Response{Err: errors.New(err.Error()), Status: 400, Result: data})
			return
		}
		claims, ok := tokenDetails.Claims.(*common.SignedDetails)
		if !ok {
			data := common.ErrorHandler("110", nil, 0, "en")
			bytes, _ := json.Marshal(data)
			(*resp).Write(bytes)
			cb(&invocation.Response{Err: errors.New("the token is invalid"), Status: 400, Result: data})
			return
		}
		if claims.ExpiresAt < time.Now().Unix() {
			data := common.ErrorHandler("111", nil, 0, "en")
			bytes, _ := json.Marshal(data)
			(*resp).Write(bytes)
			cb(&invocation.Response{Err: errors.New("token is expired"), Status: 400, Result: data})
			return
		}
	}
	chain.Next(inv, cb)
}
