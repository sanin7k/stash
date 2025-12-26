package client

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"

	"stash/internal/protocol"
)

func Backup(addr, path string) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	r := bufio.NewReader(conn)
	w := bufio.NewWriter(conn)

	fi, _ := os.Stat(path)

	fmt.Fprintf(w,
		"BACKUP\nname:%s\nsize:%d\n\n",
		filepath.Base(path), fi.Size(),
	)
	w.Flush()

	h, _ := protocol.ReadHeader(r)
	if h.Command != "OK" {
		return
	}

	offset := protocol.MustInt(h, "offset")
	size := protocol.MustInt(h, "size")

	if offset > size {
		panic("invalid offset from server")
	}

	f, _ := os.Open(path)
	defer f.Close()
	f.Seek(offset, 0)

	remaining := size - offset

	n, err := io.CopyN(w, f, remaining)
	if err != nil || n != remaining {
		panic("short upload")
	}
}

func Restore(addr, name string) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	offset := int64(0)
	if st, err := os.Stat(name); err == nil {
		offset = st.Size()
	}

	r := bufio.NewReader(conn)
	w := bufio.NewWriter(conn)

	fmt.Fprintf(w,
		"RESTORE\nname:%s\noffset:%d\n\n",
		name, offset,
	)
	w.Flush()

	h, _ := protocol.ReadHeader(r)
	if h.Command != "OK" {
		return
	}

	confirmedOffset := protocol.MustInt(h, "offset")
	size := protocol.MustInt(h, "size")

	if confirmedOffset > size {
		panic("invalid offset from server")
	}

	f, _ := os.OpenFile(name, os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	f.Seek(confirmedOffset, 0)

	remaining := size - confirmedOffset
	n, err := io.CopyN(f, r, remaining)
	if err != nil || n != remaining {
		panic("short download")
	}
}

func List(addr string) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	r := bufio.NewReader(conn)
	w := bufio.NewWriter(conn)

	w.WriteString("LIST\n\n")
	w.Flush()

	h, _ := protocol.ReadHeader(r)
	if h.Command != "OK" {
		return
	}

	count := protocol.MustInt(h, "count")
	for range count {
		line, _ := r.ReadString('\n')
		fmt.Print(line)
	}
}
