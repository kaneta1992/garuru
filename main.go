package main

import (
	"errors"
	"fmt"
	_ "io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	_ "strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/kaneta1992/simple-web-benchmarker/util"
)

func GetRandomUrl(baseUrl *url.URL, response *http.Response) (*url.URL, error) {
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

func worker(baseUrlString string, ch chan<- int, end chan bool) {
	s := session.NewSession()
	baseUrl, _ := url.Parse(baseUrlString)
	for {
		select {
		case <-end:
			return
		default:
			fmt.Printf("%v\n", baseUrl)
			req, err := s.NewRequest("GET", baseUrl.String(), nil)
			if err != nil {
				return
			}

			resp, err := s.SendRequest(req)
			if err != nil {
				return
			}
			defer resp.Body.Close()

			fmt.Printf("%d\n", resp.StatusCode)
			ch <- resp.StatusCode

			baseUrl, err = GetRandomUrl(baseUrl, resp)
			if err != nil {
				baseUrl, _ = url.Parse(baseUrlString)
				continue
			}
		}
	}
}

func counter(statusMap map[int]int, status <-chan int, end <-chan bool) {
	for {
		select {
		case v := <-status:
			statusMap[v]++
		case <-end:
			return
		}
	}
}

func main() {
	num := 2
	rand.Seed(55301)
	status := make(chan int, num)
	end := make(chan bool)

	go func(end chan<- bool) {
		time.Sleep(3 * time.Second)
		close(end)
	}(end)

	wg := new(sync.WaitGroup)
	for i := 0; i < num; i++ {
		go func() {
			wg.Add(1)
			worker("http://xn--u9j013yjqe.xn--u8jxb0b.com/", status, end)
			wg.Done()
		}()
	}

	statusMap := make(map[int]int)
	counter(statusMap, status, end)
	wg.Wait()
	fmt.Printf("%v\n", statusMap)
}
