package session

import (
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"time"

	"github.com/kaneta1992/garuru/util/cache"
)

const (
	UserAgent = "garuru"
)

type Session struct {
	Client    *http.Client
	Transport *http.Transport

	cacheStore *cache.CacheStore
	logger     *log.Logger
}

func NewSession() *Session {
	w := &Session{
		logger:     log.New(os.Stdout, "", 0),
		cacheStore: cache.NewCacheStore(),
	}

	jar, _ := cookiejar.New(&cookiejar.Options{})
	w.Transport = &http.Transport{}
	w.Client = &http.Client{
		Transport: w.Transport,
		Jar:       jar,
		Timeout:   time.Duration(10) * time.Second,
	}

	return w
}

func (s *Session) NewRequest(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	return req, err
}

func (s *Session) NewPostFormRequest(url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest("POST", url, body)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req, err
}

func (s *Session) RefreshClient() {
	jar, _ := cookiejar.New(&cookiejar.Options{})
	s.Transport = &http.Transport{}
	s.Client = &http.Client{
		Transport: s.Transport,
		Jar:       jar,
	}
}

func (s *Session) SendRequest(req *http.Request) (*http.Response, *cache.URLCache, error) {
	urlCache, cacheFound := s.cacheStore.Get(req.URL.String())
	if cacheFound {
		urlCache.Apply(req)
	} else {
		urlCache = nil
	}
	req.Header.Set("User-Agent", UserAgent)

	response, err := s.Client.Do(req)
	if err != nil {
		return nil, nil, err
	}

	if response.StatusCode == 200 {
		uc := cache.NewURLCache(response)
		if uc != nil {
			s.cacheStore.Set(req.URL.String(), uc)
		}
	}

	return response, urlCache, nil
}
