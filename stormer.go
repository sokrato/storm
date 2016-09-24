package storm

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const (
	ConsumeBufSize = 1024
)

type Stormer struct {
	cfg  Config
	stop bool
}

type Stats struct {
	Total, NFailure uint
}

func NewStormer(cfg Config) *Stormer {
	return &Stormer{cfg, false}
}

func (d *Stormer) Stop() {
	fmt.Println("Stopping...")
	d.stop = true
}

func (d *Stormer) listenForInterrupts() {
	// listen for INT to stop
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	go func() {
		cnt := 0
		for range sigCh {
			if cnt > 1 {
				os.Exit(1)
			}
			cnt++
			d.Stop()
			fmt.Printf("Press again to force quit.\n")
		}
	}()
}

func (d *Stormer) Start() {
	fmt.Println("Storm the front!")
	d.listenForInterrupts()

	startAt := time.Now()
	num := d.cfg.Concurrency()
	statsCh := make(chan Stats, 8)
	for i := 0; i < num; i++ {
		go d.runForever(statsCh)
	}

	var total, nfail uint = 0, 0
	for i := 0; i < num; i++ {
		st := <-statsCh
		total += st.Total
		nfail += st.NFailure
	}
	fmt.Printf("Total: %v, Failure: %v, Time: %v\n", total, nfail, time.Now().Sub(startAt))
}

func (d *Stormer) runForever(statsCh chan<- Stats) {
	st := Stats{0, 0}
	for !d.stop {
		total, nfailure, err := d.runOnce()
		if err != nil && err.Error() != "unexpected EOF" {
			fmt.Printf("%v\n", err)
		}
		st.Total += total
		st.NFailure += nfailure
		time.Sleep(time.Second)
	}
	statsCh <- st
}

func (d *Stormer) runOnce() (total, nfail uint, err error) {
	conn, err := d.cfg.Dial()

	if err != nil {
		return
	}

	defer func() {
		conn.Close()
	}()

	reader := bufio.NewReader(conn)
	buf := make([]byte, ConsumeBufSize)

	for !d.stop {
		total++
		data := d.cfg.Data()
		dataLen := len(data)

		for i := 0; i < dataLen; {
			nWrite, er := conn.Write(data[i:dataLen])
			if er != nil {
				nfail++
				err = er
				return
			}
			i += nWrite
		}

		var er error
		var resp *http.Response
		if TimedCall(func() {
			resp, er = http.ReadResponse(reader, nil)
		}, d.cfg.ReadTimeout()) {
			er = ErrTimeout
		}

		if er != nil {
			err = er
			nfail++
			return
		}

		keepAlive := resp.Header.Get("Connection")
		if keepAlive != "keep-alive" {
			break
		}

		for {
			n, _ := resp.Body.Read(buf)
			if n < ConsumeBufSize {
				break
			}
		}
	}
	return
}
