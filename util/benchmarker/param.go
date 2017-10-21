package benchmarker

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type Param interface {
	SetMultiPart(string, *multipart.Writer)
	SetValue(string, url.Values)
}

type TextParam string

func (p TextParam) SetMultiPart(key string, writer *multipart.Writer) {
	_ = writer.WriteField(key, string(p))
}
func (p TextParam) SetValue(key string, values url.Values) {
	values.Add(key, string(p))
}

type FileParam string

func escapeQuotes(s string) string {
	return strings.NewReplacer("\\", "\\\\", `"`, "\\\"").Replace(s)
}
func (p FileParam) SetMultiPart(key string, writer *multipart.Writer) {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			escapeQuotes(key), escapeQuotes(filepath.Base(string(p)))))
	h.Set("Content-Type", "text/plain")
	part, err := writer.CreatePart(h)
	if err != nil {
		// ここにきたらやばい
		return
	}

	file, err := os.Open(string(p))
	if err != nil {
		// ここにきたらやばい
		return
	}
	defer file.Close()

	_, err = io.Copy(part, file)
	if err != nil {
		// ここにきたらやばい
		return
	}
}
func (p FileParam) SetValue(key string, values url.Values) {
	values.Add(key, string(p))
}

type Params interface {
	NewRequest(string, string) *http.Request
	ExistKey(key string) bool
	AddKey(key string)
	GetKeys() []string
	AddParam(key string, val Param)
}

type defaultParams struct {
	params map[string][]Param
}

func (f *defaultParams) ExistKey(key string) bool {
	_, e := f.params[key]
	return e
}
func (f *defaultParams) AddKey(key string) {
	if f.ExistKey(key) {
		return
	}
	f.params[key] = []Param{}
}
func (f *defaultParams) GetKeys() []string {
	var keys []string
	for key, _ := range f.params {
		keys = append(keys, key)
	}
	return keys
}
func (f *defaultParams) AddParam(key string, val Param) {
	f.params[key] = append(f.params[key], val)
}

type EncodedParams struct {
	defaultParams
}

func NewEncodedParams() *EncodedParams {
	p := &EncodedParams{
		defaultParams: defaultParams{
			map[string][]Param{},
		},
	}
	return p
}

func (p *EncodedParams) NewRequest(method, uri string) *http.Request {
	u := url.Values{}
	for key, v := range p.params {
		for _, param := range v {
			param.SetValue(key, u)
		}
	}
	req, err := http.NewRequest(method, uri, strings.NewReader(u.Encode()))
	if err != nil {
		// ここにきたらやばい
		return nil
	}
	if method == "GET" {
		req.URL.RawQuery = u.Encode()
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	fmt.Printf("\n%v\n", req)
	return req
}

type MultiPartParams struct {
	defaultParams
}

func NewMultiPartParams() *MultiPartParams {
	p := &MultiPartParams{
		defaultParams: defaultParams{
			map[string][]Param{},
		},
	}
	return p
}

func (p *MultiPartParams) NewRequest(method, uri string) *http.Request {
	if method == "GET" {
		// ここにきたらやばい
		return nil
	}
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	for key, v := range p.params {
		for _, param := range v {
			param.SetMultiPart(key, writer)
		}
	}
	err := writer.Close()
	if err != nil {
		// ここにきたらやばい
		return nil
	}

	req, err := http.NewRequest(method, uri, body)
	if err != nil {
		// ここにきたらやばい
		return nil
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())
	fmt.Printf("\n%v\n", req)
	return req
}
