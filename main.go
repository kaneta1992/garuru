package main

import (
	"fmt"
	"github.com/kaneta1992/garuru/util/benchmarker"
)

func main() {
	num := 32

	b := benchmarker.NewBenchmarker("http://xn--u9j013yjqe.xn--u8jxb0b.com/", num)
	ret := b.Start(3)

	fmt.Printf("%v\n", ret)
}
