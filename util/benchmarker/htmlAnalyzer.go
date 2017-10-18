package benchmarker

import (
	"fmt"
	_ "io/ioutil"
	"net/http"
	"net/url"

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
			if err != nil {
				return
			}
			urls = append(urls, linkUrl)
		}
	})

	size := len(urls)
	if size < 1 {
		return nil, fmt.Errorf("not exist link")
	}
	return urls, nil
}

func (h *HttpAnalyzer) GetForms() ([]*HtmlForm, error) {
	forms := []*HtmlForm{}
	h.document.Find("form").Each(func(_ int, s *goquery.Selection) {
		form, err := NewHtmlForm(s)
		if err != nil {
			fmt.Println(err)
			return
		}
		action, _ := h.finalUrl.Parse(form.Action)
		form.Action = action.String()
		forms = append(forms, form)
	})

	size := len(forms)
	if size < 1 {
		return nil, fmt.Errorf("not exist form")
	}
	return forms, nil
}
