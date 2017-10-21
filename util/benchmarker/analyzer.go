package benchmarker

import (
	"fmt"
	"io"
	_ "io/ioutil"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

type HttpAnalyzer struct {
	finalUrl *url.URL
	document *goquery.Document
}

func NewHttpAnalyzer(finalUrl *url.URL, bodyReader io.Reader) (*HttpAnalyzer, error) {
	doc, err := goquery.NewDocumentFromReader(bodyReader)
	if err != nil {
		return nil, err
	}
	h := &HttpAnalyzer{
		finalUrl: finalUrl,
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
func (h *HttpAnalyzer) GetResourcesURL() ([]*url.URL, error) {
	urls := make([]*url.URL, 0)
	h.document.Find("img,script").Each(func(_ int, s *goquery.Selection) {
		src, exists := s.Attr("src")
		if !exists {
			return
		}
		srcUrl, err := h.finalUrl.Parse(src)
		if err != nil {
			return
		}
		urls = append(urls, srcUrl)
	})

	h.document.Find("link").Each(func(_ int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists {
			return
		}
		srcUrl, err := h.finalUrl.Parse(href)
		if err != nil {
			return
		}
		urls = append(urls, srcUrl)
	})

	size := len(urls)
	if size < 1 {
		return nil, fmt.Errorf("not exist resource")
	}
	return urls, nil
}
