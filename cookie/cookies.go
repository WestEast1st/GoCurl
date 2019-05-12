package cookie

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

//Cookies cookie関連の操作をするメソッド
type Cookies interface {
	LoadFile(string) (Cookies, error)
	WriteFile(string) error
	Read(string) ([]*http.Cookie, error)
	Add(http.Cookie) error
	Remove(domein string, name string) error
	Updata(domein string, name string, value string) error
	UpdataCookies(domein string, cookies []*http.Cookie) error
}

type cookies struct {
	data map[string][]*http.Cookie
}

// 非効率読み込み いずれどうにかする
func (c *cookies) LoadFile(filepath string) (Cookies, error) {
	var d map[string][]*http.Cookie
	d = map[string][]*http.Cookie{}

	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)

	for _, cks := range strings.Split(string(b), "\n")[4:] {
		ck := strings.Split(cks, "\t")
		if len(ck) >= 7 {
			e, _ := strconv.ParseInt(ck[4], 10, 64)
			h, _ := strconv.ParseBool(ck[1])
			s, _ := strconv.ParseBool(ck[3])
			if ck[2] == "" {
				ck[2] = "/"
			}
			d[ck[0]] = append(d[ck[0]], &http.Cookie{
				Name:       ck[5],
				Value:      ck[6],
				Path:       ck[2],
				Domain:     ck[0],
				Expires:    time.Unix(e, 0),
				RawExpires: ck[4],
				Secure:     h,
				HttpOnly:   s,
				Raw:        cks,
			})
		}
	}
	return &cookies{data: d}, nil
}

// 非効率な書き込み
func (c *cookies) WriteFile(filepath string) error {
	var jarslice []string
	for k, cks := range c.data {
		for _, ck := range cks {
			if ck.Path == "" {
				ck.Path = "/"
			}
			if ck.RawExpires == "" {
				ck.RawExpires = strconv.FormatInt(time.Now().Unix(), 10)
			}
			jarslice = append(jarslice, strings.Join([]string{
				k,
				strconv.FormatBool(ck.HttpOnly),
				ck.Path,
				strconv.FormatBool(ck.Secure),
				ck.RawExpires,
				ck.Name,
				ck.Value,
			}, "\t"))
		}
	}
	if _, err := os.Stat(path.Dir(filepath)); os.IsNotExist(err) {
		os.Mkdir(path.Dir(filepath), 0777)
	}
	file, err := os.Create(filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	head := "# Netscape HTTP Cookie File\n# http://www.netscape.com/newsref/std/cookie_spec.html\n# This is a generated file!  Do not edit.\n\n"
	file.Write(([]byte)(head + strings.Join(jarslice, "\n")))
	return nil
}

// ドメインに所属するcookie structを配列で
func (c *cookies) Read(domein string) ([]*http.Cookie, error) {
	return c.data[domein], nil
}

// cookie structを追加
func (c *cookies) Add(cookieS http.Cookie) error {
	if c.data[cookieS.Domain] == nil {
		c.data = map[string][]*http.Cookie{}
	}
	c.data[cookieS.Domain] = append(c.data[cookieS.Domain], &cookieS)
	return nil
}

//ドメイン配下のnameを削除
func (c *cookies) Remove(domein string, name string) error {
	newCookies := []*http.Cookie{}
	for _, cookie := range c.data[domein] {
		if cookie.Name != name {
			newCookies = append(newCookies, cookie)
		}
	}
	c.data[domein] = newCookies
	return nil
}

//ドメイン配下のnameの値をvalueで更新
func (c *cookies) Updata(domein string, name string, value string) error {
	for _, cookie := range c.data[domein] {
		if cookie.Name == name {
			cookie.Value = value
		}
	}
	return nil
}
func (c *cookies) UpdataCookies(domein string, cookies []*http.Cookie) error {
	if c.data[domein] == nil {
		c.data = map[string][]*http.Cookie{}
	}
	c.data[domein] = nil
	c.data[domein] = cookies
	return nil
}

//New is
func New() Cookies {
	return &cookies{}
}
