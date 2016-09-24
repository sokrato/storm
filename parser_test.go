package storm

import (
	"bytes"
	"testing"
)

func TestParseHeader(t *testing.T) {
	ua := "User-Agent: storm/1.0"
	lines := []string{
		"# This is a comment",
		ua,
	}
	header, err := ParseHeader(lines)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	if header.Get("User-Agent") != "storm/1.0" {
		t.Fail()
	}

	buf := new(bytes.Buffer)
	WriteHeader(buf, header)
	data := buf.Bytes()
	if len(data) != (len(ua) + 2) {
		t.Fail()
	}
}
