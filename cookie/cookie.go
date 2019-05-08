package cookie

import (
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

// cookie関連の操作をするメソッド
type Cookies interface {
	LoadFile(string) (Cookies, error)
	WriteFile(string) error
	Read(string) ([]cookie, error)
	Add(cookie) error
	Remove(domein string, name string) error
	Updata(domein string, name string, value string) error
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
	data map[string][]cookie
}

// 非効率読み込み いずれどうにかする
func (c *cookies) LoadFile(filepath string) (Cookies, error) {
	var d map[string][]cookie
	d = map[string][]cookie{}

	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)

	for _, cks := range strings.Split(string(b), "\n")[4:] {

		ck := strings.Split(cks, "\t")
		e, _ := strconv.ParseInt(ck[4], 10, 64)
		f, _ := strconv.ParseBool(ck[1])
		s, _ := strconv.ParseBool(ck[3])
		d[ck[0]] = append(d[ck[0]], cookie{
			ck[0],
			f,
			ck[2],
			s,
			e,
			ck[5],
			ck[6],
		})

	}
	return &cookies{data: d}, nil
}

// 非効率な書き込み
func (c *cookies) WriteFile(filepath string) error {
	return nil
}

// ドメインに所属するcookie structを配列で
func (c *cookies) Read(domein string) ([]cookie, error) {
	return c.data[domein], nil
}

// cookie structを追加
func (c *cookies) Add(cookie_s cookie) error {
	c.data[cookie_s.Domain] = append(c.data[cookie_s.Domain], cookie_s)
	return nil
}

//ドメイン配下のnameを削除
func (c *cookies) Remove(domein string, name string) error {
	return nil
}

//ドメイン配下のnameの値をvalueで更新
func (c *cookies) Updata(domein string, name string, value string) error {
	return nil
}

func New() Cookies {
	return &cookies{}
}
