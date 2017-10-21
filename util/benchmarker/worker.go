package benchmarker

import (
	"bytes"
	"fmt"
	"io"
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
	formSetter     *FormSetter
}

func NewWorker(fs *FormSetter, status chan<- int, end <-chan bool) *Worker {
	w := &Worker{
		httpSession:    session.NewSession(),
		responseStatus: status,
		endBroadCaster: end,
		formSetter:     fs,
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

func getRandomForm(analyzer *HttpAnalyzer) (*HtmlForm, error) {
	forms, err := analyzer.GetForms()
	if err != nil {
		return nil, err
	}
	return forms[rand.Intn(len(forms))], nil
}

func (w *Worker) createRequestFromResponse(analyzer *HttpAnalyzer) (*http.Request, error) {
	method := rand.Intn(2)
	for i := 0; i < 2; i++ {
		switch method {
		case 0: //GET
			nextUrl, err := getRandomUrl(analyzer)
			if err != nil {
				method = 1
				continue
			}
			return w.httpSession.NewRequest("GET", nextUrl.String(), nil)
		case 1: //POST
			form, err := getRandomForm(analyzer)
			if err != nil {
				method = 0
				continue
			}
			w.formSetter.Set(form)
			return form.BuildRequest(), nil
		}
	}
	return nil, fmt.Errorf("not exists url")
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
		res, _, err := w.httpSession.SendRequest(req)
		if err != nil {
			fmt.Printf("error resourse send request: %s\n", v.String())
			continue
		}
		fmt.Printf("%d\n", res.StatusCode)
		select {
		case <-w.endBroadCaster:
		case w.responseStatus <- res.StatusCode:
		}
		res.Body.Close()
	}
	return nil
}

func (w *Worker) Start(startUrl string) error {
	request, err := w.httpSession.NewRequest("GET", startUrl, nil)
	if err != nil {
		return err
	}
	for {
		select {
		case <-w.endBroadCaster:
			return nil
		default:
			fmt.Printf("%v\n", request.URL.String())

			response, cache, err := w.httpSession.SendRequest(request)
			if err != nil || response.Body == nil {
				fmt.Printf("error request: %s\n", err)
				request, _ = w.httpSession.NewRequest("GET", startUrl, nil)
				continue
			}

			var body io.Reader = response.Body
			if response.StatusCode == 304 {
				body = bytes.NewReader(cache.Body)
			}

			fmt.Printf("%d\n", response.StatusCode)
			select {
			case <-w.endBroadCaster:
			case w.responseStatus <- response.StatusCode:
			}

			analyzer, err := NewHttpAnalyzer(response.Request.URL, body)
			w.getResources(analyzer)
			request, err = w.createRequestFromResponse(analyzer)
			if err != nil {
				fmt.Printf("create request error: %v\n", err)
				request, _ = w.httpSession.NewRequest("GET", startUrl, nil)
			}
		}
	}
}
