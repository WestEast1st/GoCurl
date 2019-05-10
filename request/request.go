package request

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

// Request is http.Request & client interfase
type Request interface {
	SetHeader([]string) error
	SetCookie([]*http.Cookie) error
	UpdataMethod(string) error
	UpdataData(string) error
	UpdataURL(string) error
	Do() (*http.Response, error)
	Read() (request, error)
	Load(request) error
}

type request struct {
	Method     string
	URL        string
	Data       string
	Headers    map[string]string
	Cookie     []*http.Cookie
	IsRedirect bool
}

func (r *request) Read() (request, error) {
	return request{
		Method:     r.Method,
		URL:        r.URL,
		Data:       r.Data,
		Headers:    r.Headers,
		Cookie:     r.Cookie,
		IsRedirect: r.IsRedirect,
	}, nil
}
func (r *request) Load(in request) error {
	r = &in
	return nil
}
func (r *request) SetHeader(headers []string) error {
	for _, header := range headers {
		headerSlice := strings.Split(header, ":")
		r.Headers[headerSlice[0]] = strings.Join(headerSlice[1:], ":")
	}
	return nil
}
func (r *request) SetCookie(cookies []*http.Cookie) error {
	r.Cookie = cookies
	return nil
}

func (r *request) UpdataMethod(method string) error {
	r.Method = method
	return nil
}

func (r *request) UpdataData(data string) error {
	r.Data = data
	return nil
}

func (r *request) UpdataURL(url string) error {
	r.URL = url
	return nil
}

func (r *request) Do() (*http.Response, error) {
	jar, _ := cookiejar.New(nil)
	req, _ := http.NewRequest(r.Method, r.URL, strings.NewReader(r.Data))
	for k, v := range r.Headers {
		req.Header.Add(k, v)
	}
	u, _ := url.Parse(r.URL)
	jar.SetCookies(u, r.Cookie)
	httpclient := http.Client{Jar: jar}
	if r.IsRedirect {
		httpclient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}
	resp, _ := httpclient.Do(req)
	return resp, nil
}

// New is Reqest interfase return
func New(method string, url string, data string) Request {
	return &request{
		Method:     method,
		URL:        url,
		Data:       data,
		Headers:    map[string]string{},
		Cookie:     []*http.Cookie{},
		IsRedirect: false,
	}
}
