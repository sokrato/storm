package dos

import (
	"crypto/tls"
	"io"
	"net"
)

type Conn interface {
	io.ReadWriteCloser
}

type Dialer interface {
	Dial() (Conn, error)
}

type CommonDialer struct {
	addr string
}

func (d CommonDialer) Dial() (Conn, error) {
	return net.Dial("tcp", d.addr)
}

type TlsDialer struct {
	addr      string
	tlsConfig *tls.Config
}

func (d TlsDialer) Dial() (Conn, error) {
	return tls.Dial("tcp", d.addr, d.tlsConfig)
}
