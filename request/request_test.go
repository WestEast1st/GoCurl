package request

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	URL    = "http://example.com"
	METHOD = "GET"
	DATA   = "hoge=fuga&piyo=hogehoge"
)

var (
	headConst = []string{
		"accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3",
		"accept-language: ja",
		"cache-control: max-age=0",
	}
	hogeCookie = http.Cookie{
		Name:   "hoge",
		Value:  "fuga",
		Path:   "/",
		Domain: "example.com",
	}
	trueCookie = http.Cookie{
		Name:   "ok",
		Value:  "true",
		Path:   "/",
		Domain: "example.com",
	}
	cookieConst = []*http.Cookie{
		&hogeCookie,
		&trueCookie,
	}
	reqConst = request{
		Method:     METHOD,
		URL:        URL,
		Data:       DATA,
		Headers:    map[string]string{},
		Cookie:     []*http.Cookie{},
		IsRedirect: false,
	}
)
var testReq = New(METHOD, URL, "")

func TestRead(t *testing.T) {
	req, _ := testReq.Read()
	if req.Data == "hoge=fuga&piyo=hogehoge" {
		t.Error("No new instance has been created")
	}
}
func TestLoada(t *testing.T) {
	testReq.Load(reqConst)
	req, _ := testReq.Read()
	if req.Headers == nil {
		t.Error("No new instance has been created")
	}
}
func TestSetHeader(t *testing.T) {
	if err := testReq.SetHeader(headConst); err != nil {
		t.Error(err)
	}
	reqst, _ := testReq.Read()
	if reqst.Headers["accept"] != " text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3" {
		t.Error("header is not set.")
	}
}

func TestSetCookie(t *testing.T) {
	if err := testReq.SetCookie(cookieConst); err != nil {
		t.Error(err)
	}
	reqst, _ := testReq.Read()
	if reqst.Cookie[1].Name != "ok" {
		t.Error("Cookie is not set.")
	}
}
func TestDo(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Send response to be tested
		rw.Write([]byte(`OK`))
	}))
	// Close the server when test finishes
	defer server.Close()
	testReq = New(METHOD, server.URL, "")
	res, _ := testReq.Do()
	defer res.Body.Close()
	if res.StatusCode != 200 {
		t.Error("not connect")
	}

}
