package dos

import (
	"log"
	"testing"
)

func TestStormer(t *testing.T) {
	cfg, _ := NewSimpleConfig("GET", "http://localhost/msg.json", 1, 2, []byte("Host: localhost\nUser-Agent: dripper/1.0\n\n"))
	dripper := NewStormer(*cfg)
	log.Println(dripper)
	dripper.Start()
}
