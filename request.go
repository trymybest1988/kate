package kate

import (
	"net/http"
	"net/url"

	"github.com/julienschmidt/httprouter"
)

type Request struct {
	*http.Request

	Id        uint64
	Env       map[string]interface{}
	RoutePath string
	Path      httprouter.Params
	RawBody   []byte
}

func (r *Request) BaseUrl() *url.URL {
	scheme := r.URL.Scheme
	if scheme == "" {
		scheme = "http"
	}

	host := r.Host
	if len(host) > 0 && host[len(host)-1] == '/' {
		host = host[:len(host)-1]
	}

	return &url.URL{
		Scheme: scheme,
		Host:   host,
	}
}

func (r *Request) UrlFor(path string, queryParams map[string][]string) *url.URL {
	baseUrl := r.BaseUrl()
	baseUrl.Path = path
	if queryParams != nil {
		query := url.Values{}
		for k, v := range queryParams {
			for _, vv := range v {
				query.Add(k, vv)
			}
		}
		baseUrl.RawQuery = query.Encode()
	}
	return baseUrl
}
