/**
 * Sample Chassis Handler to print log
 *
**/

package chassisHandlers

import (
	"VoD_microservice/common"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/emicklei/go-restful"
	"github.com/go-chassis/go-archaius"
	"github.com/go-chassis/go-chassis/v2/core/handler"
	"github.com/go-chassis/go-chassis/v2/core/invocation"
	"github.com/go-chassis/openlog"
	"github.com/xeipuuv/gojsonschema"
)

const Name = "Payload-Validator"

type PayloadValidatorHanldlerHandler struct{}

func init() { handler.RegisterHandler(Name, New) }

func New() handler.Handler { return &PayloadValidatorHanldlerHandler{} }

func (h *PayloadValidatorHanldlerHandler) Name() string { return Name }

func (h *PayloadValidatorHanldlerHandler) Handle(chain *handler.Chain, inv *invocation.Invocation, cb invocation.ResponseCallBack) {
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

	payload_bytes, err := ioutil.ReadAll(req.Body)
	// Restore the io.ReadCloser to its original state
	req.Body = ioutil.NopCloser(bytes.NewBuffer(payload_bytes))
	openlog.Debug("got request to " + inv.URLPath)
	schemaPath := getSchema(inv.OperationID)
	openlog.Info(inv.OperationID)
	openlog.Info(schemaPath)
	if schemaPath == "" {
		chain.Next(inv, func(r *invocation.Response) {
			cb(r)
		})
		return
	}

	schema_bytes, _ := ioutil.ReadFile(schemaPath)
	schemaLoader := gojsonschema.NewBytesLoader(schema_bytes)
	documentLoader := gojsonschema.NewBytesLoader(payload_bytes)
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		openlog.Error("error occured here" + err.Error())
		data := common.ErrorHandler("201", result, 0, "en")
		bytes, _ := json.Marshal(data)
		(*resp).WriteHeader(http.StatusBadRequest)
		(*resp).Write(bytes)
		cb(&invocation.Response{Err: errors.New(data.Msg), Status: 400, Result: data})
		return
	}
	if result.Valid() {
		openlog.Info("Payload Validation completed")
		chain.Next(inv, cb)
		return
	} else {

		validationErrors := make([]string, 0)
		for _, desc := range result.Errors() {
			validationErrors = append(validationErrors, desc.String())
		}
		fmt.Println(validationErrors, "######################################")
		data := common.ErrorHandler("100", validationErrors, 0, "en")
		bytes, _ := json.Marshal(data)
		(*resp).WriteHeader(http.StatusBadRequest)
		(*resp).Write(bytes)
		cb(&invocation.Response{Err: errors.New(data.Msg), Status: 400, Result: data})
		openlog.Error("schema validation errors")
		return
	}
}

func getSchema(api string) string {
	return archaius.GetString(api, "")
}
