package storm

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	entitySep = []byte("\r\n\r\n")
)

type Config interface {
	Concurrency() int
	RequestsPerThread() int
	ReadTimeout() time.Duration
	Data() []byte
	Show()
	Dialer
}

type SimpleConfig struct {
	address           string
	concurrency       int
	requestsPerThread int
	readTimeout       time.Duration
	data              []byte
	dialer            Dialer
	// secure           bool
	// tlsConfig        *tls.Config
}

func (c SimpleConfig) Concurrency() int {
	return c.concurrency
}
func (c SimpleConfig) RequestsPerThread() int {
	return c.requestsPerThread
}
func (c SimpleConfig) ReadTimeout() time.Duration {
	return c.readTimeout
}
func (c SimpleConfig) Data() []byte {
	return c.data
}
func (c SimpleConfig) Show() {
	fmt.Printf("Address: %v\n", c.address)
	fmt.Printf("Concurrency: %v\n", c.concurrency)
	fmt.Printf("RequestsPerThread: %v\n", c.requestsPerThread)
	fmt.Printf("data: \n---------\n%v\n", string(c.data))
}
func (c SimpleConfig) Dial() (Conn, error) {
	return c.dialer.Dial()
}

func parseUrl(rawurl string) (addr, path string, secure bool, err error) {
	aurl, err := url.Parse(rawurl)
	if err != nil {
		return
	}

	addr = aurl.Host
	secure = aurl.Scheme == "https"
	if strings.Index(addr, ":") < 0 {
		if secure {
			addr = addr + ":443"
		} else {
			addr = addr + ":80"
		}
	}
	path = aurl.Path
	if aurl.RawQuery != "" {
		path += "?" + aurl.RawQuery
	}
	if aurl.Fragment != "" {
		path += "#" + aurl.Fragment
	}
	return
}

func NewSimpleConfig(method, rawurl string, concurrency, readTimeout int, header http.Header, entity []byte) (*SimpleConfig, error) {
	addr, path, secure, err := parseUrl(rawurl)
	if err != nil {
		return nil, err
	}

	var dialer Dialer = nil
	if secure {
		dialer = TlsDialer{addr, &tls.Config{}}
	} else {
		dialer = CommonDialer{addr}
	}

	statusLine := []byte(method + " " + path + " HTTP/1.1" + LineSep)
	buf := bytes.NewBuffer(statusLine)

	header.Set("Connection", "keep-alive")
	header.Del("Content-Length")
	WriteHeader(buf, header)

	buf.Write(LineSepBytes)
	buf.Write(LineSepBytes)

	entityLength := len(entity)
	if entityLength > 0 {
		header.Set("Content-Length", fmt.Sprintf("%v", entityLength))
		buf.Write(entity)
	}

	return &SimpleConfig{
		address:           addr,
		concurrency:       concurrency,
		requestsPerThread: 0,
		readTimeout:       time.Duration(readTimeout) * time.Millisecond,
		data:              buf.Bytes(),
		dialer:            dialer,
	}, nil
}
