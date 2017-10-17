package benchmarker

import (
	"fmt"
	_ "io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	_ "strings"

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
	analyzer, err := NewHttpAnalyzer(response)
	urls, err := analyzer.GetLinks()
	if err != nil {
		return nil, err
	}
	return urls[rand.Intn(len(urls))], nil
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
