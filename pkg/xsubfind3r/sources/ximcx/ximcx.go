package ximcx

import (
	"encoding/json"
	"fmt"

	"github.com/hueristiq/xsubfind3r/pkg/xsubfind3r/httpclient"
	"github.com/hueristiq/xsubfind3r/pkg/xsubfind3r/sources"
)

type Source struct{}

type response struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
	Data    []struct {
		Domain string `json:"domain"`
	} `json:"data"`
}

func (source *Source) Run(config *sources.Configuration) (subdomains chan sources.Subdomain) {
	subdomains = make(chan sources.Subdomain)

	go func() {
		defer close(subdomains)

		res, err := httpclient.SimpleGet(fmt.Sprintf("http://sbd.ximcx.cn/DomainServlet?domain=%s", config.Domain))
		if err != nil {
			return
		}

		var results response

		if err := json.Unmarshal(res.Body(), &results); err != nil {
			return
		}

		for _, result := range results.Data {
			subdomains <- sources.Subdomain{Source: source.Name(), Value: result.Domain}
		}
	}()

	return
}

func (source *Source) Name() string {
	return "ximcx"
}