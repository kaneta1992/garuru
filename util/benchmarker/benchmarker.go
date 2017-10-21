package benchmarker

import (
	_ "io/ioutil"
	"log"
	_ "strings"
	"sync"
	"time"
)

type Benchmarker struct {
	StartUrl string

	worker         []*Worker
	statusCounter  map[int]int
	responseStatus chan int
	endBroadCaster chan bool
}

func NewBenchmarker(startUrl string, configPath string, workerNum int) *Benchmarker {
	b := &Benchmarker{
		StartUrl:       startUrl,
		statusCounter:  make(map[int]int),
		responseStatus: make(chan int, workerNum),
		endBroadCaster: make(chan bool),
	}

	var formSetter *FormSetter = nil
	var err error
	if configPath != "" {
		formSetter, err = NewFormSetter(configPath)
		if err != nil {
			log.Fatalf("テストデータ生成に失敗")
		}
	}
	for i := 0; i < workerNum; i++ {
		w := NewWorker(formSetter, b.responseStatus, b.endBroadCaster)
		b.worker = append(b.worker, w)
	}
	return b
}

func (b *Benchmarker) Start(second time.Duration) map[int]int {
	go func(end chan<- bool) {
		time.Sleep(second * time.Second)
		close(end)
	}(b.endBroadCaster)

	wg := new(sync.WaitGroup)
	for _, w := range b.worker {
		go func(worker *Worker) {
			wg.Add(1)
			worker.Start(b.StartUrl)
			wg.Done()
		}(w)
	}

	for {
		select {
		case status := <-b.responseStatus:
			b.statusCounter[status]++
		case <-b.endBroadCaster:
			wg.Wait()
			return b.statusCounter
		}
	}
}
