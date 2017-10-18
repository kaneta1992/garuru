package benchmarker

import (
	"fmt"
	_ "io/ioutil"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type HtmlForm struct {
	Action  string
	Method  string
	EncType string
	Params
}

func (f *HtmlForm) attrValidate() error {
	if f.Action == "" {
		return fmt.Errorf("not exists action: %s", f.Action)
	}
	if f.Method == "" {
		return fmt.Errorf("read method error: %s", f.Method)
	}
	switch f.EncType {
	case "multipart/form-data":
		f.Params = NewMultiPartParams()
	case "application/x-www-form-urlencoded":
		f.Params = NewEncodedParams()
	default:
		return fmt.Errorf("unknown enctype: %s", f.EncType)
	}
	return nil
}

func NewHtmlForm(s *goquery.Selection) (*HtmlForm, error) {
	f := &HtmlForm{}
	var exists bool
	f.Action, _ = s.Attr("action")
	f.Method, _ = s.Attr("method")
	f.EncType, exists = s.Attr("enctype")
	// enctypeはデフォルト値がある
	if !exists {
		f.EncType = "application/x-www-form-urlencoded"
	}
	err := f.attrValidate()
	if err != nil {
		return nil, err
	}
	s.Find("input,select,textarea").Each(func(_ int, param *goquery.Selection) {
		name, exists := param.Attr("name")
		if !exists {
			return
		}
		f.AddKey(name)
		t, _ := param.Attr("type")
		switch t {
		case "hidden":
			// Attrはkeyが見つからなかった場合に空文字が帰ってくる
			v, _ := param.Attr("value")
			f.AddParam(name, TextParam(v))
		}
	})
	return f, nil
}

func (f *HtmlForm) BuildRequest() *http.Request {
	return f.Params.NewRequest(strings.ToUpper(f.Method), f.Action)
}
