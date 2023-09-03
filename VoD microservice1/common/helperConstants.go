package common

var DbName string

var CollectionName = "video"

var SECRET_KEY string

var Errorcodes = make(map[string]interface{})

type HTTPResponse struct {
	Msg        string      `json:"_msg"`
	Status     int         `json:"_status"`
	Data       interface{} `json:"_data"`
	TotalCount int64       `json:"_totalcount"`
	ErrorCode  int64       `json:"_errorcode"`
}
