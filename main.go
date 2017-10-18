package main

import (
	"fmt"
	"github.com/kaneta1992/garuru/util/benchmarker"
	"github.com/kaneta1992/garuru/util/session"
	"io/ioutil"
)

func main() {
	// num := 32

	// b := benchmarker.NewBenchmarker("http://xn--u9j013yjqe.xn--u8jxb0b.com/", num)
	// ret := b.Start(3)

	// fmt.Printf("%v\n", ret)
	// s := session.NewSession()
	// req, _ := s.NewRequest("GET", "http://sed.jp/uploader.php", nil)
	// res, _ := s.SendRequest(req)
	// h, _ := benchmarker.NewHttpAnalyzer(res)
	// forms, _ := h.GetForms()
	// forms[2].AddKey("upfile")
	// forms[2].AddParam("upfile", benchmarker.FileParam("./main.go"))
	// forms[2].AddKey("t")
	// forms[2].AddParam("t", benchmarker.TextParam("yaxa"))

	// req = forms[2].BuildRequest()
	// res, _ = s.SendRequest(req)

	// b, _ := ioutil.ReadAll(res.Body)
	// fmt.Printf("%d\n", res.StatusCode)
	// fmt.Printf("%v\n", string(b))
	// res.Body.Close()
	s := session.NewSession()
	p := benchmarker.NewMultiPartParams()
	p.AddParam("upfile", benchmarker.FileParam("./config.yml"))
	p.AddParam("t", benchmarker.TextParam("textだよ"))
	req := p.NewRequest("POST", "http://sed.jp/upload.php")
	res, _ := s.SendRequest(req)
	b, _ := ioutil.ReadAll(res.Body)
	fmt.Printf("%d\n", res.StatusCode)
	fmt.Printf("%v\n", string(b))
	res.Body.Close()
}
