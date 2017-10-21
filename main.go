package main

import (
	"fmt"
	"github.com/kaneta1992/garuru/util/benchmarker"
	_ "github.com/kaneta1992/garuru/util/session"
	_ "io/ioutil"
)

func main() {
	num := 32

	b := benchmarker.NewBenchmarker("http://sed.jp/uploader.html", "test.yml", num)
	ret := b.Start(3)
	fmt.Printf("%v\n", ret)
	return
}
