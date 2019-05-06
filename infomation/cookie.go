package infomation

import (
	"os"
	"path"
	"strconv"
	"strings"
)

// cookie関連の操作をするメソッド
type Cookies interface {
	LoadCookieFile(string) (Cookies, error)
	WriteCookieFile(string) error
	ReadCookie() (string, error)
	AddCookie(setdata string) error
	RemoveCookie(name string) error
	UpdataCookie(name string, value string) error
}

// Cookie.txt出力用
// Netscape HTTP Cookie File
//http://www.netscape.com/newsref/std/cookie_spec.html
// This is a generated file!  Do not edit.
type cookie struct {
	Domain     string
	Flag       bool
	Path       string
	Secure     bool
	Expiration int64 //有効期限
	Name       string
	Value      string
}

type cookies struct {
	Data map[string]cookie
}

// 非効率読み込み いずれどうにかする
func (c *cookies) LoadCookieFile(filepath string) (Cookies, error) {
	var d map[string]cookie

	return &cookies{Data: d}, nil
}

// 非効率書き込み いずれどうにかする
func (c *cookies) WriteCookieFile(filepath string) error {
	var jarslice []string
	for _, ck := range c.Data {
		jarslice = append(jarslice, strings.Join([]string{ck.Domain,
			strconv.FormatBool(ck.Flag),
			ck.Path,
			strconv.FormatBool(ck.Secure),
			strconv.FormatInt(ck.Expiration, 10),
			ck.Name,
			ck.Value,
		}, "\t"))
	}
	if _, err := os.Stat(path.Dir(filepath)); os.IsNotExist(err) {
		os.Mkdir(path.Dir(filepath), 0777)
	}
	file, err := os.Create(filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	file.Write(([]byte)(strings.Join(jarslice, "\n")))
	return nil
}

func (c *cookies) ReadCookie() (string, error) {
	return "", nil
}

func (c *cookies) AddCookie(setdata string) error {
	return nil
}

func (c *cookies) RemoveCookie(name string) error {
	return nil
}

func (c *cookies) UpdataCookie(name string, value string) error {
	return nil
}

func New() Cookies {
	return &cookies{}
}
