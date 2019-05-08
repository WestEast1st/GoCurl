package cookie

import (
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"testing"
	"time"
)

var randSrc = rand.NewSource(time.Now().UnixNano())

const (
	rsLetters       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	rsLetterIdxBits = 6
	rsLetterIdxMask = 1<<rsLetterIdxBits - 1
	rsLetterIdxMax  = 63 / rsLetterIdxBits
)

func RandString(n int) string {
	b := make([]byte, n)
	cache, remain := randSrc.Int63(), rsLetterIdxMax
	for i := n - 1; i >= 0; {
		if remain == 0 {
			cache, remain = randSrc.Int63(), rsLetterIdxMax
		}
		idx := int(cache & rsLetterIdxMask)
		if idx < len(rsLetters) {
			b[i] = rsLetters[idx]
			i--
		}
		cache >>= rsLetterIdxBits
		remain--
	}
	return string(b)
}

var testfilepath = "./test_cookiejar.txt"
var c = New()

func TestNew(t *testing.T) {
	reflect.TypeOf(c)
}

//cookie読み込みテスト
func TestReadCookie(t *testing.T) {
	cookies, _ := c.LoadFile(testfilepath)
	a, _ := cookies.Read("home.netscape.com")
	if a[0].Domain != "home.netscape.com" {
		t.Error("取得した値が異なります\n")
	}
}

//cookiejar読み込みテスト
func TestLoadCookieFile(t *testing.T) {
	cookies, err := c.LoadFile(testfilepath)
	if err != nil {
		t.Errorf("error message :%v", err)
	}
	a, _ := cookies.Read("home.netscape.com")
	if a[0].Domain != "home.netscape.com" {
		t.Error("取得した値が異なります\n")
	}
	b, _ := cookies.Read("www2n.meshnet.or.jp")
	if b[3].Name != "fb04" {
		t.Error("取得した値が異なります\n")
	}
	c, _ := cookies.Read(".netscape.com")
	if c[1].Secure != false {
		t.Error("取得した値が異なります\n")
	}
}

var test_cookies, _ = c.LoadFile(testfilepath)

//cookie追加のテスト
func TestAdd(t *testing.T) {
	rs := RandString(10)
	c := cookie{
		"www.geocities.com",
		false,
		"/",
		false,
		934947883,
		"LocalId",
		rs,
	}
	e := test_cookies.Add(c)
	if e != nil {
		t.Errorf("error message :%v", e)
	}
	a, _ := test_cookies.Read("www.geocities.com")
	if a[len(a)-1].Value != rs {
		t.Error("クッキーの追加がなされていません\n")
	}
}

//cookie更新
func TestUpdata(t *testing.T) {
	rs := RandString(10)
	e := test_cookies.Updata("www.geocities.com", "LocalId", rs)
	if e != nil {
		t.Errorf("error message :%v", e)
	}
	a, _ := test_cookies.Read("www.geocities.com")
	if a[len(a)-1].Value != rs {
		t.Error("クッキーの更新がなされていません\n")
	}
}

//cookie削除のテスト
func TestRemove(t *testing.T) {
	e := test_cookies.Remove("www.geocities.com", "LocalId")
	if e != nil {
		t.Errorf("error message :%v", e)
	}
	a, _ := test_cookies.Read("www.geocities.com")
	if len(a) != 1 {
		t.Error("クッキーの削除がなされていません\n")
	}
}

//cookiejar書き込みのテスト
func TestWriteFile(t *testing.T) {
	test_file := "./test_cookiejar_test.txt"
	test_cookies.WriteFile(test_file)
	c, _ = c.LoadFile(test_file)
	test_slice, _ := c.Read("www.enemy.org")
	if len(test_slice) != 1 {
		t.Error("書き込まれた内容にロスが確認できます\n")
	}
	if err := os.Remove(test_file); err != nil {
		fmt.Println(err)
	}
}
