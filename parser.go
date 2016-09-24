package storm

import (
	"bufio"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var (
	HeaderSep      = ": "
	HeaderSepBytes = []byte(HeaderSep)
	LineSep        = "\r\n"
	LineSepBytes   = []byte(LineSep)
)

func ParseHeader(lines []string) (header http.Header, err error) {
	header = make(http.Header)
	for _, line := range lines {
		if strings.HasPrefix(line, "#") {
			continue
		}
		kv := strings.SplitN(line, HeaderSep, 2)
		if len(kv) != 2 {
			err = errors.New("Bad Header Line: " + line)
			return
		}
		name := strings.TrimSpace(kv[0])
		value := strings.TrimSpace(kv[1])
		header.Add(name, value)
	}
	return
}

func WriteHeader(output io.Writer, header http.Header) error {
	for name, values := range header {
		_, err := output.Write([]byte(name))
		if err != nil {
			return err
		}

		output.Write(HeaderSepBytes)

		for idx, value := range values {
			if idx > 0 {
				output.Write([]byte("; "))
			}
			output.Write([]byte(value))
		}
		output.Write(LineSepBytes)
	}
	return nil
}

func ParseHeaderAndEntity(input *bufio.Reader) (header http.Header, entity []byte, err error) {
	var line string
	headerLines := make([]string, 0, 8)

	for {
		line, err = input.ReadString('\n')
		if err != nil && err != io.EOF {
			return
		}
		line = strings.TrimSpace(line)
		if line == "" { // end of header lines
			break
		}
		headerLines = append(headerLines, line)
	}

	header, err = ParseHeader(headerLines)
	if err != nil {
		return
	}

	entity, err = ioutil.ReadAll(input)
	if err == io.EOF {
		err = nil
	}

	if err != nil {
		return
	}

	return
}

func ReadHeaderAndEntityFromFile(filename string) (http.Header, []byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}

	reader := bufio.NewReader(file)
	return ParseHeaderAndEntity(reader)
}
