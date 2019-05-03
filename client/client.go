package client

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/urfave/cli"
)

//Client ...
//clientはhttpリクエストを送信するためのモジュールです。
type Client interface {
	Get() (string, error)
	GetImage()
}

type client struct {
	URL      string
	Method   string
	FileName string
	Data     string
}

func (c *client) GetImage() {
	response, err := http.Get(c.URL)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	if _, err := os.Stat(path.Dir(c.FileName)); os.IsNotExist(err) {
		os.Mkdir(path.Dir(c.FileName), 0777)
	}
	file, err := os.Create(c.FileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	io.Copy(file, response.Body)
}

func (c *client) Get() (string, error) {
	req, _ := http.NewRequest("GET", c.URL, nil)

	httpclient := new(http.Client)
	resp, err := httpclient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)
	return string(byteArray), nil
}

func getUrl(c *cli.Context) string {
	var urls []string
	r := regexp.MustCompile(`^(http|https)://([\w-]+\.)+[\w-]+(/[\w-./?%&=]*)?$`)
	for _, arg := range c.Args() {
		if len(r.FindStringSubmatch(arg)) > 0 {
			urls = append(urls, r.FindStringSubmatch(arg)[0])
		}
	}
	return urls[len(urls)-1]
}

func New(c *cli.Context) Client {
	var method string
	var urlstring string
	var filename string

	urlstring = getUrl(c)

	if c.Bool("O") {
		u, err := url.Parse(urlstring)
		if err != nil {
			log.Fatal(err)
		}
		if u.Path == "/" {
			filename = "index.html"
		} else {
			filename = u.Path
		}
	} else if c.String("output") != "" {
		filename = c.String("output")
	} else {
		filename = ""
	}

	if c.String("data") != "" {
		method = "POST"
	} else {
		method = "GET"
	}

	return &client{
		URL:      urlstring,
		Method:   method,
		FileName: "." + filename,
		Data:     strings.Join(c.StringSlice("data"), "&"),
	}
}
