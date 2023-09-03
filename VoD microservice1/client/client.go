package client

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-chassis/foundation/httputil"
	"github.com/go-chassis/go-chassis/core"
	"github.com/go-chassis/go-chassis/v2/client/rest"
	"github.com/go-chassis/openlog"
)

func GetClient(method string, url string, headers map[string]string, payload []byte) (map[string]interface{}, error) {

	req, err := rest.NewRequest(url, method, payload)
	if err != nil {
		openlog.Error(err.Error())
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := core.NewRestInvoker().ContextDo(context.TODO(), req)
	if err != nil {
		openlog.Error("do request failed. : " + err.Error())
		return nil, err
	}
	body := httputil.ReadBody(resp)
	fmt.Println("body: ", string(body))
	var res = make(map[string]interface{})
	err = json.Unmarshal(body, &res)
	if err != nil {
		openlog.Error(err.Error())
		return res, err
	}
	return res, nil

}
