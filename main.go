package main

import (
	"fmt"
	"github.com/go-yaml/yaml"
	_ "github.com/kaneta1992/garuru/util/benchmarker"
	_ "github.com/kaneta1992/garuru/util/session"
	"io/ioutil"
	_ "io/ioutil"
)

type YmlParams struct {
	Action  string
	Enctype string
	Method  string
	Data    []struct {
		Types  map[string]string
		Values []map[string][]string
	}
}

func main() {
	buf, err := ioutil.ReadFile("config.yml")
	if err != nil {
		panic(err)
	}

	var d []YmlParams
	err = yaml.Unmarshal(buf, &d)
	fmt.Println(d[0].Data[0].Types["one"])
	fmt.Println(d[0].Data[0].Types["two"])
	fmt.Println(d[0].Data[0].Values[0]["one"])
	fmt.Println(d[0].Data[0].Values[0]["two"])
	fmt.Println(d[0].Data[0].Values[1]["one"])
	fmt.Println(d[0].Data[0].Values[1]["two"])
	return

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

	// s := session.NewSession()
	// p := benchmarker.NewMultiPartParams()
	// p.AddParam("upfile", benchmarker.FileParam("./config.yml"))
	// p.AddParam("t", benchmarker.TextParam("textだよ"))
	// req := p.NewRequest("POST", "http://sed.jp/upload.php")
	// res, _ := s.SendRequest(req)
	// b, _ := ioutil.ReadAll(res.Body)
	// fmt.Printf("%d\n", res.StatusCode)
	// fmt.Printf("%v\n", string(b))
	// res.Body.Close()
}
