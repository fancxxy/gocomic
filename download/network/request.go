package network

import (
	"compress/gzip"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var client = &http.Client{
	Timeout: 1 * time.Minute,
}

// Request contains request client
type Request struct {
	client *http.Client
}

// Header represents http request header
type Header map[string]string

// Query represents http request query param
type Query = url.Values

func (r *Request) do(method, rawurl string, vs ...interface{}) (*Response, error) {
	if rawurl == "" {
		return nil, errors.New("request.do: url is not specified")
	}

	var (
		res *http.Response
		req = &http.Request{
			Method:     method,
			Header:     make(http.Header),
			Proto:      "HTTP/1.1",
			ProtoMajor: 1,
			ProtoMinor: 1,
		}
		query, form Query
		response    = &Response{Req: req}
	)

	for _, v := range vs {
		switch vv := v.(type) {
		case Header:
			for key, value := range vv {
				req.Header.Add(key, value)
			}
		case Query:
			if method == "GET" {
				query = vv
			} else {
				form = vv
			}
		}
	}

	if form != nil {
		body := form.Encode()
		response.reqBody = []byte(body)
		req.Body = ioutil.NopCloser(strings.NewReader(body))
		if req.Header.Get("Content-Type") == "" {
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
		}
	}

	if query != nil {
		rawurl = rawurl + "?" + query.Encode()
	}
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	req.URL = u

	if host := req.Header.Get("Host"); host != "" {
		req.Host = host
	}

	res, err = r.client.Do(req)
	if err != nil {
		return nil, err
	}

	response.Res = res

	if res.Header.Get("Content-Encoding") == "gzip" && req.Header.Get("Accept-Encoding") != "" {
		body, err := gzip.NewReader(res.Body)
		if err != nil {
			return nil, err
		}
		res.Body = body
	}

	return response, nil
}

// New create *Request
func New() *Request {
	req := &Request{client: client}
	return req
}

// Get request
func (r *Request) Get(url string, v ...interface{}) (*Response, error) {
	return r.do("GET", url, v...)
}

// Post request
func (r *Request) Post(url string, v ...interface{}) (*Response, error) {
	return r.do("POST", url, v...)
}
