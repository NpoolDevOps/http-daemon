package httpdaemon

import (
	"fmt"
	log "github.com/EntropyPool/entropy-logger"
	"github.com/go-resty/resty/v2"
)

var cli = resty.New()

func R() *resty.Request {
	return cli.R()
}

func ParseResponse(resp *resty.Response) (*ApiResp, error) {
	apiResp, err := ParseResponseBody(resp.Body())
	if err != nil {
		log.Errorf(log.Fields{}, "< api response parse error %v [%v] (%v)",
			resp.Request.URL, resp.Request.QueryParam, err)
		return nil, err
	}

	if apiResp.Code != 0 {
		log.Errorf(log.Fields{}, "< api response error %v [%v] (%v)",
			resp.Request.URL, resp.Request.QueryParam, apiResp.Msg)
		return nil, fmt.Errorf("%v (%v)", apiResp.Msg, apiResp.Code)
	}

	return apiResp, nil
}
