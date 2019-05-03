package infomation

import "net/url"

// 送信するHTTP関連の情報
type HttpInfomation struct {
	URL      string
	URI      string
	Method   string
	Port     string
	Hostname string
	Path     string
	Query    url.Values
	Data     url.Values
	Fragment string
	Output   Output
	Header   Header
}

// アウトプットの関連情報
type Output struct {
	Flag     bool
	Filepath string //ユーザーが配置したいディレクトリ
	Filename string //命名のPATH
}

type Header struct {
	ReadFlag bool
}
