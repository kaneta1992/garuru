package cache

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/marcw/cachecontrol"
)

type CacheStore struct {
	items map[string]*URLCache
}

func NewCacheStore() *CacheStore {
	m := make(map[string]*URLCache)
	c := &CacheStore{
		items: m,
	}
	return c
}

func (c *CacheStore) Get(key string) (*URLCache, bool) {
	v, found := c.items[key]
	return v, found
}

func (c *CacheStore) Set(key string, value *URLCache) {
	c.items[key] = value
}

type URLCache struct {
	LastModified string
	Etag         string
	ExpiresAt    time.Time
	CacheControl *cachecontrol.CacheControl
	Body         []byte
}

func NewURLCache(res *http.Response) *URLCache {
	directive := res.Header.Get("Cache-Control")
	cc := cachecontrol.Parse(directive)
	noCache, _ := cc.NoCache()

	if len(directive) == 0 || noCache || cc.NoStore() {
		return nil
	}

	now := time.Now()
	lm := res.Header.Get("Last-Modified")
	etag := res.Header.Get("ETag")

	b, _ := ioutil.ReadAll(res.Body)
	res.Body = ioutil.NopCloser(bytes.NewReader(b))

	return &URLCache{
		LastModified: lm,
		Etag:         etag,
		ExpiresAt:    now.Add(cc.MaxAge()),
		CacheControl: &cc,
		Body:         b,
	}
}

func (c *URLCache) Available() bool {
	return time.Now().Before(c.ExpiresAt)
}

func (c *URLCache) Apply(req *http.Request) {
	if c.Available() {
		if c.LastModified != "" {
			req.Header.Add("If-Modified-Since", c.LastModified)
		}

		if c.Etag != "" {
			req.Header.Add("If-None-Match", c.Etag)
		}
	}
}
