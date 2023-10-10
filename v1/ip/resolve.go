package ip

import (
	"encoding/json"
	"github.com/parnurzeal/gorequest"
)

type ResolveResp struct {
	Code    int
	Data    data
	Message string
}

type data struct {
	CountryCode   string
	CountryName   string
	ContinentCode string
	CityName      string
	Timezone      string
}

// Resolve 解析海外IP
func Resolve(ip string) (ResolveResp, []error) {
	var resp ResolveResp
	var errs []error

	_, body, errs := gorequest.New().
		Get("http://ip.gtarcade.com/api/resolve").
		Query("ip=" + ip).
		End()
	if len(errs) > 0 {
		resp.Code = 500
		return resp, errs
	}

	err := json.Unmarshal([]byte(body), &resp)
	if err != nil {
		resp.Code = 500
		errs = append(errs, err)
		return resp, errs
	}

	return resp, nil
}
