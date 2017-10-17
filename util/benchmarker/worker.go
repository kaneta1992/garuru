package benchmarker

import (
	"errors"
	"fmt"
	_ "io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	_ "strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/kaneta1992/simple-web-benchmarker/util/session"
)

type Worker struct {
	httpSession    *session.Session
	responseStatus chan<- int
	endBroadCaster <-chan bool
}

func NewWorker(status chan<- int, end <-chan bool) *Worker {
	w := &Worker{
		httpSession:    session.NewSession(),
		responseStatus: status,
		endBroadCaster: end,
	}
	return w
}

func getRandomUrl(response *http.Response) (*url.URL, error) {
	baseUrl := response.Request.URL
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return nil, err
	}

	urls := make([]string, 0)
	doc.Find("a").Each(func(_ int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			reqUrl, err := baseUrl.Parse(href)
			if err == nil {
				urls = append(urls, reqUrl.String())
			}
		}
	})

	size := len(urls)
	if size < 1 {
		return nil, errors.New("not exist link")
	}

	return url.Parse(urls[rand.Intn(size)])
}

func (w *Worker) createRequestFromResponse(response *http.Response) (*http.Request, error) {
	nextUrl, err := getRandomUrl(response)
	if err != nil {
		return nil, err
	}
	return w.httpSession.NewRequest("GET", nextUrl.String(), nil)
}

func (w *Worker) Start(startUrl string) error {
	defaultRequest, err := w.httpSession.NewRequest("GET", startUrl, nil)
	request := defaultRequest
	if err != nil {
		return err
	}
	for {
		select {
		case <-w.endBroadCaster:
			return nil
		default:
			fmt.Printf("%v\n", request.URL.String())
			response, err := w.httpSession.SendRequest(request)
			if err != nil {
				continue
			}

			fmt.Printf("%d\n", response.StatusCode)
			w.responseStatus <- response.StatusCode

			request, err = w.createRequestFromResponse(response)
			if err != nil {
				request = defaultRequest
			}
		}
	}
}
