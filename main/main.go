package main

import (
	"flag"
	"github.com/dlutxx/storm"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

var (
	url               = flag.String("url", "http://localhost/msg.json", "url to storm")
	method            = flag.String("method", "GET", "http verb")
	concurrency       = flag.Int("concurrency", 64, "how many concurrent connections")
	requestsPerThread = 0 // flag.Int("requestsPerThread", 0, "how many requests each thread send")
	datafile          = flag.String("datafile", "data.txt", "which file to read")
	timeToRun         = flag.Int("ttr", 0, "time to run, seconds")
	test              = flag.Bool("test", false, "show config and exit")
)

func init() {
	flag.Parse()
	*method = strings.ToUpper(*method)
}

func main() {
	data, err := ioutil.ReadFile(*datafile)
	if err != nil {
		log.Fatalln(err)
	}
	cfg, _ := storm.NewSimpleConfig(*method, *url, *concurrency, requestsPerThread, data)
	if *test {
		cfg.Show()
		return
	}

	stormer := storm.NewStormer(*cfg)

	if *timeToRun > 0 {
		go func() {
			time.Sleep(time.Duration(*timeToRun) * time.Second)
			stormer.Stop()
		}()
	}

	stormer.Start()
}
