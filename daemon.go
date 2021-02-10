package httpdaemon

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

type HttpHandler func(w http.ResponseWriter, req *http.Request) (interface{}, string, int)

type HttpRouter struct {
	Location string
	Handler  HttpHandler
	Method   string
}

var routerTable []HttpRouter = make([]HttpRouter, 0)

type ApiResp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Body interface{} `json:"body"`
}

func response(w http.ResponseWriter, resp interface{}, msg string, code int) error {
	apiResp := ApiResp{
		Code: code,
		Msg:  msg,
		Body: resp,
	}
	jsonStr, err := json.Marshal(&apiResp)
	if nil != err {
		return err
	}
	w.Write(jsonStr)
	return nil
}

func rootHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("request %v -> %v [%v]", req.RemoteAddr, req.URL, req.URL.Path)
	if err := req.ParseForm(); nil != err {
		log.Printf("fail to parse form %v", req.URL)
		response(w, struct{}{}, err.Error(), -1)
		return
	}
	for _, r := range routerTable {
		if r.Location == req.URL.Path {
			if r.Method != req.Method {
				log.Printf("request method is %v, need %v", req.Method, r.Method)
				return
			}
			resp, msg, code := r.Handler(w, req)
			err := response(w, resp, msg, code)
			if nil != err {
				log.Printf("fail to response %v", req.URL)
			}

			return
		}
	}

	response(w, struct{}{}, fmt.Sprintf("unknown request %v", req.URL), -3)
}

func Run(port int) error {
	http.HandleFunc("/", rootHandler)

	go func(port int) {
		portStr := fmt.Sprintf(":%v", port)
		log.Printf("start http server [%v]", portStr)
		for {
			http.ListenAndServe(portStr, nil)
		}
	}(port)

	return nil
}

func RegisterRouter(router HttpRouter) error {
	for _, r := range routerTable {
		if r.Location == router.Location {
			return errors.New("router already exist")
		}
	}
	log.Printf("add router: %v", router.Location)
	routerTable = append(routerTable, router)
	return nil
}

func ValidateParams(keys []string, params map[string][]string) error {
	var err error
	for _, key := range keys {
		if len(params[key]) == 0 || params[key][0] == "" {
			err = errors.New("params are not matched or empty")
			break
		}
	}

	return err
}

func ParseResponseBody(resBody []byte) *ApiResp {
	var unmar map[string]interface{}
	_ = json.Unmarshal(resBody, &unmar)

	parseRes := new(ApiResp)

	parseRes.Code = int(unmar["code"].(float64))
	if _, ok := unmar["msg"]; ok {
		parseRes.Msg = unmar["msg"].(string)
	} else {
		parseRes.Msg = unmar["error"].(string)
	}
	parseRes.Body = unmar["body"]

	return parseRes
}
