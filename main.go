package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"

	"github.com/ivahaev/go-logger"
	"github.com/kaneta1992/garuru/util/benchmarker"
)

func main() {
	var (
		verbose    bool
		num        int
		configPath string
		t          int
	)
	flag.IntVar(&t, "t", 3, "ベンチマークの時間")
	flag.IntVar(&t, "time", 3, "")
	flag.BoolVar(&verbose, "v", false, "ログレベル")
	flag.BoolVar(&verbose, "verbose", false, "")
	flag.IntVar(&num, "w", 4, "ワーカー数")
	flag.IntVar(&num, "worker", 4, "")
	flag.StringVar(&configPath, "p", "", "POSTのテストデータ設定ファイルのパス")
	flag.StringVar(&configPath, "post", "", "")
	flag.Parse()

	logger.SetLevel("CRIT")
	if verbose {
		logger.SetLevel("INFO")
	}

	rand.Seed(time.Now().UnixNano())
	fmt.Printf("%v\n", flag.Args())
	b := benchmarker.NewBenchmarker(flag.Args(), configPath, num)
	ret := b.Start(time.Duration(t))
	fmt.Printf("%v\n", ret)
	return
}
