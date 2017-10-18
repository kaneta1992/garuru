package session

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	UserAgent = "garuru"
)

type Session struct {
	Client    *http.Client
	Transport *http.Transport

	logger *log.Logger
}

func NewSession() *Session {
	w := &Session{
		logger: log.New(os.Stdout, "", 0),
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

func (s *Session) SendRequest(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", UserAgent)

	return s.Client.Do(req)
}
