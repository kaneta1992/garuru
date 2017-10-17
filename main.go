package main

import (
	"github.com/kaneta1992/simple-web-benchmarker/util/benchmarker"
	"github.com/kaneta1992/simple-web-benchmarker/util/session"
	"github.com/tonnerre/golang-pretty"
)

func main() {
	// num := 32

	// b := benchmarker.NewBenchmarker("http://xn--u9j013yjqe.xn--u8jxb0b.com/", num)
	// ret := b.Start(3)

	// fmt.Printf("%v\n", ret)
	s := session.NewSession()
	req, _ := s.NewRequest("GET", "http://sed.jp/uploader.php", nil)
	res, _ := s.SendRequest(req)
	h, _ := benchmarker.NewHttpAnalyzer(res)
	forms, _ := h.GetForms()
	pretty.Printf("%v\n", *forms[0])
}
