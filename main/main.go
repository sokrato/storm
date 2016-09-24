package main

import (
	"flag"
	"github.com/dlutxx/storm"
	"log"
	"strings"
	"time"
)

var (
	url               = flag.String("url", "http://localhost/", "url to storm")
	method            = flag.String("method", "GET", "http verb")
	concurrency       = flag.Int("concurrency", 64, "how many concurrent connections")
	requestsPerThread = 0 //flag.Int("requestsPerThread", 0, "how many requests each thread send") // not supported yet
	readTimeout       = flag.Int("readTimeout", 200, "response read timeout in millisecond")
	requestData       = flag.String("requestData", "request.txt", "which file to read")
	timeToRun         = flag.Int("ttr", 0, "time to run, seconds")
	show              = flag.Bool("show", false, "show config and exit")
)

func init() {
	flag.Parse()
	*method = strings.ToUpper(*method)
}

func main() {
	header, entity, err := storm.ReadHeaderAndEntityFromFile(*requestData)
	if err != nil {
		log.Fatalln("Failed to read request data: " + err.Error())
	}

	*method = strings.ToUpper(*method)
	cfg, _ := storm.NewSimpleConfig(*method, *url, *concurrency, *readTimeout, header, entity)
	if *show {
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
