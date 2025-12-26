package protocol

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

type Header struct {
	Command string
	Fields  map[string]string
}

func ReadHeader(r *bufio.Reader) (*Header, error) {
	line, err := r.ReadString('\n')
	if err != nil {
		return nil, err
	}

	h := &Header{
		Command: strings.TrimSpace(line),
		Fields:  make(map[string]string),
	}

	for {
		l, err := r.ReadString('\n')
		if err != nil {
			return nil, err
		}
		l = strings.TrimSpace(l)
		if l == "" {
			break
		}
		parts := strings.SplitN(l, ":", 2)
		if len(parts) != 2 {
			return nil, ErrMalformedHeader
		}
		h.Fields[parts[0]] = parts[1]
	}

	return h, nil
}

func WriteOK(w *bufio.Writer, fields map[string]int64) {
	w.WriteString("OK\n")
	for k, v := range fields {
		fmt.Fprintf(w, "%s:%d\n", k, v)
	}
	w.WriteString("\n")
}

func WriteError(w *bufio.Writer, msg string) {
	fmt.Fprintf(w, "ERR\nmsg:%s\n\n", msg)
}

func MustInt(h *Header, key string) int64 {
	v, ok := h.Fields[key]
	if !ok {
		panic("missing field: " + key)
	}
	n, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		panic("invalid int field: " + key)
	}
	return n
}

func ReadOffset(r *bufio.Reader) int64 {
	h, err := ReadHeader(r)
	if err != nil || h.Command != "OK" {
		return 0
	}
	o, _ := strconv.ParseInt(h.Fields["offset"], 10, 64)
	return o
}
