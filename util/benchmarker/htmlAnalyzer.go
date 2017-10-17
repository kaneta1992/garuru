package benchmarker

import (
	"errors"
	_ "io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type HttpAnalyzer struct {
	finalUrl *url.URL
	document *goquery.Document
}

func NewHttpAnalyzer(response *http.Response) (*HttpAnalyzer, error) {
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return nil, err
	}
	h := &HttpAnalyzer{
		finalUrl: response.Request.URL,
		document: doc,
	}
	return h, err
}

func (h *HttpAnalyzer) GetLinks() ([]*url.URL, error) {
	urls := []*url.URL{}
	h.document.Find("a").Each(func(_ int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			linkUrl, err := h.finalUrl.Parse(href)
			if err == nil {
				urls = append(urls, linkUrl)
			}
		}
	})

	size := len(urls)
	if size < 1 {
		return nil, errors.New("not exist link")
	}
	return urls, nil
}

type HtmlForm struct {
	Action  string
	Method  string
	EncType string
	Params  map[string][]string
}

func NewHtmlForm(action, method, enctype string) *HtmlForm {
	h := &HtmlForm{
		Action:  action,
		Method:  method,
		EncType: enctype,
		Params:  make(map[string][]string),
	}
	return h
}

func (h *HttpAnalyzer) GetForms() ([]*HtmlForm, error) {
	forms := []*HtmlForm{}
	h.document.Find("form").Each(func(_ int, fs *goquery.Selection) {
		action, exists := fs.Attr("action")
		if !exists {
			return
		}
		method, exists := fs.Attr("method")
		// POST以外のformは考慮していない
		if !exists || !strings.EqualFold(method, "post") {
			return
		}
		enctype, exists := fs.Attr("enctype")
		if !exists {
			enctype = "application/x-www-form-urlencoded"
		}
		form := NewHtmlForm(action, method, enctype)
		fs.Find("input,select,textarea").Each(func(_ int, s *goquery.Selection) {
			name, exists := s.Attr("name")
			if exists {
				if _, exists = form.Params[name]; !exists {
					form.Params[name] = []string{}
				}
				t, _ := s.Attr("type")
				switch t {
				case "hidden":
					if v, exists := s.Attr("value"); exists {
						form.Params[name] = append(form.Params[name], v)
					}
				}
			}
		})
		forms = append(forms, form)
	})

	size := len(forms)
	if size < 1 {
		return nil, errors.New("not exist form")
	}
	return forms, nil
}
