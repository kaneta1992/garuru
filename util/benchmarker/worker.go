package benchmarker

import (
	"fmt"
	_ "io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	_ "strings"

	"github.com/kaneta1992/garuru/util/session"
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

func getRandomUrl(analyzer *HttpAnalyzer) (*url.URL, error) {
	urls, err := analyzer.GetLinks()
	if err != nil {
		return nil, err
	}
	return urls[rand.Intn(len(urls))], nil
}

func (w *Worker) createRequestFromResponse(analyzer *HttpAnalyzer) (*http.Request, error) {
	nextUrl, err := getRandomUrl(analyzer)
	if err != nil {
		return nil, err
	}
	return w.httpSession.NewRequest("GET", nextUrl.String(), nil)
}

func (w *Worker) getResources(analyzer *HttpAnalyzer) error {
	urls, err := analyzer.GetResourcesURL()
	if err != nil {
		return err
	}
	for _, v := range urls {
		fmt.Printf("%v\n", v.String())
		req, err := w.httpSession.NewRequest("GET", v.String(), nil)
		if err != nil {
			fmt.Printf("error resourse new request: %s\n", v.String())
			continue
		}
		res, err := w.httpSession.SendRequest(req)
		if err != nil {
			fmt.Printf("error resourse send request: %s\n", v.String())
			continue
		}
		fmt.Printf("%d\n", res.StatusCode)
		res.Body.Close()
	}
	return nil
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

			analyzer, err := NewHttpAnalyzer(response)
			w.getResources(analyzer)
			request, err = w.createRequestFromResponse(analyzer)
			if err != nil {
				request = defaultRequest
			}
		}
	}
}
