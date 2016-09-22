package storm

import (
	"crypto/tls"
	"testing"
)

func TestCommonDialer(t *testing.T) {
	dialer := CommonDialer{"baidu.com:80"}
	conn, err := dialer.Dial()
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	conn.Close()
}

func TestTlsDialer(t *testing.T) {
	dialer := TlsDialer{"baidu.com:443", &tls.Config{}}
	conn, err := dialer.Dial()
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	conn.Close()
}
