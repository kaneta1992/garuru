package benchmarker

import (
	"errors"
	_ "io/ioutil"
	"net/http"
	"net/url"
	_ "strings"

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
	urls := make([]*url.URL, 0)
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
		return nil, errors.New("not exist resource")
	}
	return urls, nil
}
